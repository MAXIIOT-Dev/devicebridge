/*
 * @Description: lora数据处理服务
 * @Copyright: Maxiiot(c) 2019
 * @Author: tgq
 * @LastEditors: tgq
 * @Date: 2019-04-11 17:00:05
 * @LastEditTime: 2019-04-19 10:09:18
 */

package server

import (
	"fmt"
	"sync"

	"github.com/maxiiot/vbaseBridge/backend/http"
	"github.com/maxiiot/vbaseBridge/backend/mqtt"
	"github.com/maxiiot/vbaseBridge/config"
	"github.com/maxiiot/vbaseBridge/storage"

	log "github.com/sirupsen/logrus"
)

var Serv *Server

// BackendServer define backend server
type BackendServer interface {
	// HandleUplinks handle uplinks
	HandleUplinks(*sync.WaitGroup)
	// Close close resource
	Close() error
	// Notice
	Notice(map[string]bool)
}

// Server 后台服务
// support http/mqtt backend
type Server struct {
	backends []BackendServer
	wg       sync.WaitGroup
}

// NewServer return server point
func NewServer(cfg config.Configuration) (*Server, error) {
	backends := make([]BackendServer, 0, 2)

	devs, err := storage.GetDevicesEUI()
	if err != nil {
		return nil, err
	}
	backends = append(backends, mqtt.NewBackend(cfg.LoraBackend.Mqtt, devs))

	var addr string
	if cfg.LoraBackend.HTTPPort > 0 {
		addr = fmt.Sprintf(":%d", cfg.LoraBackend.HTTPPort)
	} else {
		addr = ":8080"
	}
	httpserv := http.New(addr)

	backends = append(backends, httpserv)

	serv := &Server{backends: backends}

	return serv, nil
}

// Start server start
func (s *Server) Start() {
	for _, backend := range s.backends {
		go func() {
			s.wg.Add(1)
			defer s.wg.Done()
			backend.HandleUplinks(&s.wg)
		}()
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
