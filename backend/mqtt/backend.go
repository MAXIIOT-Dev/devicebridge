/*
 * @Description: mqtt mqtt消息订阅处理
 * @Copyright: Maxiiot(c) 2019
 * @Author: tgq
 * @LastEditors: tgq
 * @Date: 2019-04-11 16:57:48
 * @LastEditTime: 2019-04-24 20:14:22
 */

package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/maxiiot/devicebridge/backend"

	paho "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// Config mqtt broker configuration
type Config struct {
	Server              string `mapstructure:"server" json:"server"`
	Username            string `mapstructure:"username" json:"username"`
	Password            string `mapstructure:"password" json:"password"`
	QOS                 uint8  `mapstructure:"qos" json:"qos"`
	CleanSession        bool   `mapstructure:"clean_session" json:"clean_session"`
	ClientID            string `mapstructure:"client_id" json:"client_id"`
	CACert              string `mapstructure:"ca_cert" json:"ca_cert"`
	TLSCert             string `mapstructure:"tls_cert" json:"tls_cert"`
	TLSKey              string `mapstructure:"tls_key" json:"tls_key"`
	UplinkTopicTemplate string `mapstructure:"uplink_topic_template" json:"uplink_topic_template"`
	// AckTopicTemplate    string `mapstructure:"ack_topic_template" json:"ack_topic_template"`
}

// Backend mqtt backend
type Backend struct {
	wg           sync.WaitGroup
	config       Config
	rxPacketChan chan backend.DataUpPayloadChan
	conn         paho.Client
	subDevs      []string // 只订阅添加的设备消息
	subNotice    chan map[string]bool
}

// NewBackend retruns backend
func NewBackend(c Config, devs []string) *Backend {
	var err error
	b := Backend{
		config:       c,
		rxPacketChan: make(chan backend.DataUpPayloadChan, 10),
		subDevs:      devs,
		subNotice:    make(chan map[string]bool),
	}

	opts := paho.NewClientOptions()
	opts.AddBroker(b.config.Server)
	opts.SetUsername(b.config.Username)
	opts.SetPassword(b.config.Password)
	opts.SetCleanSession(b.config.CleanSession)
	opts.SetClientID(b.config.ClientID)
	opts.SetOnConnectHandler(b.onConnected)
	opts.SetConnectionLostHandler(b.onConnectionLost)

	tlsconfig, err := newTLSConfig(b.config.CACert, b.config.TLSCert, b.config.TLSKey)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"ca_cert":  b.config.CACert,
			"tls_cert": b.config.TLSCert,
			"tls_key":  b.config.TLSKey,
		}).Fatal("backend/mqtt: error loading mqtt certificate files")
	}

	if tlsconfig != nil {
		opts.SetTLSConfig(tlsconfig)
	}
	log.WithField("server", b.config.Server).Info("backend/mqtt: connecting to mqtt broker")
	b.conn = paho.NewClient(opts)
	for {
		if token := b.conn.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("backend/mqtt: connecting to mqtt broker failed, will retry in 2s: %s", token.Error())
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	go b.noticeHandler()

	return &b
}

func (b *Backend) onConnected(c paho.Client) {
	log.Info("backend/mqtt: connected to mqtt server")

	if len(b.subDevs) > 0 {
		topics := make(map[string]byte)
		topicsField := make([]string, 0, len(b.subDevs))
		for _, dev := range b.subDevs {
			topic := fmt.Sprintf(b.config.UplinkTopicTemplate, dev)
			topicsField = append(topicsField, topic)
			topics[topic] = b.config.QOS
		}
		for {
			log.WithFields(log.Fields{
				"topics": topicsField,
				"qos":    b.config.QOS,
			}).Info("backend/mqtt: subscribing to rx topic")
			if token := b.conn.SubscribeMultiple(topics, b.rxPacketHandler); token.Wait() && token.Error() != nil {
				log.WithFields(log.Fields{
					"topics": topicsField,
					"qos":    b.config.QOS,
				}).Errorf("backend/mqtt: subscribe error: %s", token.Error())
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
}

func (b *Backend) rxPacketHandler(c paho.Client, msg paho.Message) {
	b.wg.Add(1)
	defer b.wg.Done()

	var rxdata backend.DataUpPayload
	if err := json.Unmarshal(msg.Payload(), &rxdata); err != nil {
		log.Errorf("backend/mqtt: decode rx packet error: %s\n", err)
		return
	}

	log.WithFields(log.Fields{
		"payload": rxdata.Data,
		"device":  rxdata.DevEUI,
	}).Info("mqtt: uplink frame received")
	if data, err := hex.DecodeString(rxdata.Data); err == nil {
		dataChan := backend.DataUpPayloadChan{
			Data:   data,
			DevEUI: rxdata.DevEUI,
		}
		b.rxPacketChan <- dataChan
	} else {
		log.WithError(err).Error("hex deocde payload data error ")
	}
}

func (b *Backend) onConnectionLost(c paho.Client, reason error) {
	log.Errorf("backend/mqtt: mqtt connection error: %s", reason)
}

// Close unsubscribe message and close chans
func (b *Backend) Close() error {
	log.Info("backend/mqtt: closing backend")

	topic := fmt.Sprintf(b.config.UplinkTopicTemplate, "+")
	log.WithField("topic", topic).Info("mqtt: unsubscribing from uplink ")
	if token := b.conn.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return fmt.Errorf("backend/mqtt: unsubscribe from %s error: %s", topic, token.Error())
	}

	log.Info("backend/mqtt: handling last messages")
	b.wg.Wait()
	close(b.rxPacketChan)
	close(b.subNotice)
	return nil
}

// Notice notice
func (b *Backend) Notice(notice map[string]bool) {
	b.subNotice <- notice
}

// noticeHandler添加设备时添加消息订阅，删除设备时取消消息订阅
func (b *Backend) noticeHandler() {
	for notice := range b.subNotice {
		for dev, issub := range notice {
			if issub {
				topic := fmt.Sprintf(b.config.UplinkTopicTemplate, dev)
				log.WithFields(log.Fields{
					"topic": topic,
					"qos":   b.config.QOS,
				}).Info("backend/mqtt: subscribing to rx topic")

				if token := b.conn.Subscribe(topic, b.config.QOS, b.rxPacketHandler); token.Wait() && token.Error() != nil {
					log.WithFields(log.Fields{
						"topics": topic,
						"qos":    b.config.QOS,
					}).Errorf("backend/mqtt: subscribe error: %s", token.Error())
				}
			} else {
				topic := fmt.Sprintf(b.config.UplinkTopicTemplate, dev)
				log.WithFields(log.Fields{
					"topic": topic,
					"qos":   b.config.QOS,
				}).Info("backend/mqtt: unsubscribing to rx topic")

				if token := b.conn.Unsubscribe(topic); token.Wait() && token.Error() != nil {
					log.WithFields(log.Fields{
						"topics": topic,
						"qos":    b.config.QOS,
					}).Errorf("backend/mqtt: unsubscribe error: %s", token.Error())
				}
			}
		}
	}
}

func newTLSConfig(cafile, certFile, certKeyFile string) (*tls.Config, error) {
	if cafile == "" && certFile == "" && certKeyFile == "" {
		return nil, nil
	}

	tlsConfig := &tls.Config{}

	// Import trusted certificates from CAfile.pem.
	if cafile != "" {
		cacert, err := ioutil.ReadFile(cafile)
		if err != nil {
			log.WithError(err).Error("gateway/mqtt: could not load ca certificate")
			return nil, err
		}
		certpool := x509.NewCertPool()
		certpool.AppendCertsFromPEM(cacert)

		tlsConfig.RootCAs = certpool // RootCAs = certs used to verify server cert.
	}

	// Import certificate and the key
	if certFile != "" && certKeyFile != "" {
		kp, err := tls.LoadX509KeyPair(certFile, certKeyFile)
		if err != nil {
			log.WithError(err).Error("gateway/mqtt: could not load mqtt tls key-pair")
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{kp}
	}

	return tlsConfig, nil
}

// HandleUplinks 处理lora上行数据
func (b *Backend) HandleUplinks(conn paho.Client, wg *sync.WaitGroup) {
	// log.Println("debug: mqtt handlerup")
	for uplink := range b.rxPacketChan {
		go func(uplink backend.DataUpPayloadChan) {
			wg.Add(1)
			defer wg.Done()
			log.Println("debug:", "handle uplink.")
			if err := backend.HandleUplink(conn, uplink); err != nil {
				log.WithFields(log.Fields{
					"device": uplink.DevEUI,
					"data":   hex.EncodeToString(uplink.Data),
				}).Errorf("process device uplink data error: %s", err)
			}
		}(uplink)
	}

}
