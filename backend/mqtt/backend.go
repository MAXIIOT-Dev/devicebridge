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
	rxPacketChan chan DataUpPayloadChan
	// ackPacketChan chan ACKNotificationChan
	conn paho.Client
}

// NewBackend retruns backend
func NewBackend(c Config) (*Backend, error) {
	var err error
	b := Backend{
		config:       c,
		rxPacketChan: make(chan DataUpPayloadChan),
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

	return &b, nil
}

func (b *Backend) onConnected(c paho.Client) {
	log.Info("backend/mqtt: connected to mqtt server")

	for {
		log.WithFields(log.Fields{
			"topic": b.config.UplinkTopicTemplate,
			"qos":   b.config.QOS,
		}).Info("backend/mqtt: subscribing to rx topic")
		if token := b.conn.Subscribe(b.config.UplinkTopicTemplate, b.config.QOS, b.rxPacketHandler); token.Wait() && token.Error() != nil {
			log.WithFields(log.Fields{
				"topic": b.config.UplinkTopicTemplate,
				"qos":   b.config.QOS,
			}).Errorf("backend/mqtt: subscribe error: %s", token.Error())
			time.Sleep(time.Second)
			continue
		}
		break
	}

}

func (b *Backend) rxPacketHandler(c paho.Client, msg paho.Message) {
	b.wg.Add(1)
	defer b.wg.Done()

	log.Info("mqtt: uplink frame received")
	var rxdata DataUpPayload
	if err := json.Unmarshal(msg.Payload(), &rxdata); err != nil {
		log.Errorf("backend/mqtt: decode rx packet error: %s\n", err)
		return
	}
	//test data
	// rxdata.Data = "ff015c2c8cda16002a3f01ff"

	if data, err := hex.DecodeString(rxdata.Data); err == nil {
		dataChan := DataUpPayloadChan{
			Data:   data,
			DevEUI: rxdata.DevEUI,
		}
		b.rxPacketChan <- dataChan
	} else {
		log.WithError(err).Error("deocde payload data error ")
	}
}

func (b *Backend) onConnectionLost(c paho.Client, reason error) {
	log.Errorf("backend/mqtt: mqtt connection error: %s", reason)
}

func (b *Backend) Close() error {
	log.Info("backend/mqtt: closing backend")

	log.WithField("topic", b.config.UplinkTopicTemplate).Info("mqtt: unsubscribing from uplink ")
	if token := b.conn.Unsubscribe(b.config.UplinkTopicTemplate); token.Wait() && token.Error() != nil {
		return fmt.Errorf("backend/mqtt: unsubscribe from %s error: %s", b.config.UplinkTopicTemplate, token.Error())
	}

	log.Info("backend/mqtt: handling last messages")
	b.wg.Wait()
	close(b.rxPacketChan)
	return nil
}

func (b *Backend) RXPacketChan() chan DataUpPayloadChan {
	return b.rxPacketChan
}

// func (b *Backend) ACKPacketChan() chan ACKNotificationChan {
// 	return b.ackPacketChan
// }

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
