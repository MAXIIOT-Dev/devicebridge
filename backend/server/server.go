/*
 * @Description: lora数据处理服务
 * @Copyright: Maxiiot(c) 2019
 * @Author: tgq
 * @LastEditors: tgq
 * @Date: 2019-04-11 17:00:05
 * @LastEditTime: 2019-04-25 10:52:21
 */

package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/maxiiot/devicebridge/backend/http"
	"github.com/maxiiot/devicebridge/backend/mqtt"
	"github.com/maxiiot/devicebridge/config"
	"github.com/maxiiot/devicebridge/storage"

	paho "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var Serv *Server

// BackendServer define backend server
type BackendServer interface {
	// HandleUplinks handle uplinks
	HandleUplinks(paho.Client, *sync.WaitGroup)
	// Close close resource
	Close() error
	// Notice
	Notice(map[string]bool)
}

// Server 后台服务
// support http/mqtt backend
type Server struct {
	backends  map[string]BackendServer
	wg        sync.WaitGroup
	publisher paho.Client
}

// NewServer return server point
func NewServer(cfg config.Configuration) (*Server, error) {
	backends := make(map[string]BackendServer)

	devs, err := storage.GetDevicesEUI()
	if err != nil {
		return nil, err
	}
	backends["mqtt"] = mqtt.NewBackend(cfg.LoraBackend.Mqtt, devs)

	var addr string
	if cfg.LoraBackend.HTTPPort > 0 {
		addr = fmt.Sprintf(":%d", cfg.LoraBackend.HTTPPort)
	} else {
		addr = ":8080"
	}
	httpserv := http.New(addr)

	backends["http"] = httpserv

	conn, err := newPublisher(cfg.Publisher.Mqtt)
	if err != nil {
		return nil, err
	}
	serv := &Server{backends: backends, publisher: conn}

	return serv, nil
}

// Start server start
func (s *Server) Start() {
	for key, _ := range s.backends {
		go func(backend BackendServer) {
			s.wg.Add(1)
			defer s.wg.Done()
			backend.HandleUplinks(s.publisher, &s.wg)
		}(s.backends[key])
	}

}

// Stop 关闭相关资源
func (s *Server) Stop() error {
	for _, backend := range s.backends {
		if err := backend.Close(); err != nil {
			return err
		}
	}
	log.Info("waiting for pending actions to complete")
	s.wg.Wait()
	return nil
}

// OnDeviceChange notice backend subscribe/unsubscribe message
func (s *Server) OnDeviceChange(notice map[string]bool) {
	for _, backend := range s.backends {
		backend.Notice(notice)
	}
}

func newPublisher(cfg mqtt.Config) (paho.Client, error) {

	opts := paho.NewClientOptions()
	opts.AddBroker(cfg.Server)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)
	opts.SetCleanSession(cfg.CleanSession)
	opts.SetClientID(cfg.ClientID)
	opts.SetConnectionLostHandler(onConnectionLost)

	tlsconfig, err := newTLSConfig(cfg.CACert, cfg.TLSCert, cfg.TLSKey)
	if err != nil {
		return nil, err
	}

	if tlsconfig != nil {
		opts.SetTLSConfig(tlsconfig)
	}
	log.WithField("server", cfg.Server).Info("publisher/mqtt: connecting to mqtt broker")
	conn := paho.NewClient(opts)
	for {
		if token := conn.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("publisher/mqtt: connecting to mqtt broker failed, will retry in 2s: %s", token.Error())
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	return conn, nil
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

func onConnectionLost(c paho.Client, reason error) {
	log.Errorf("backend/mqtt: mqtt connection error: %s", reason)
}
