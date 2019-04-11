/*
 * @Description: lora数据处理服务
 * @Copyright: Maxiiot(c) 2019
 * @Author: tgq
 * @LastEditors: tgq
 * @Date: 2019-04-11 17:00:05
 * @LastEditTime: 2019-04-11 18:28:19
 */

package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/maxiiot/vbaseBridge/backend"
	"github.com/maxiiot/vbaseBridge/backend/http"
	"github.com/maxiiot/vbaseBridge/backend/mqtt"
	"github.com/maxiiot/vbaseBridge/backend/protocol"
	"github.com/maxiiot/vbaseBridge/config"
	"github.com/maxiiot/vbaseBridge/storage"

	log "github.com/sirupsen/logrus"
)

const (
	// TransportMQTT transport by mqtt
	TransportMQTT = "mqtt"
	// TransportHTTP transport by http
	TransportHTTP = "http"
)

var Serv *Server

// BackendServer define backend server
type BackendServer interface {
	// RXPacketChan return data uppayload
	RXPacketChan() chan backend.DataUpPayloadChan
	// Close close resource
	Close() error
	// Notice
	Notice(map[string]bool)
}

// Server 后台服务
type Server struct {
	backend BackendServer
	wg      sync.WaitGroup
}

// NewServer return server point
func NewServer(cfg config.Configuration) (*Server, error) {
	var typ string
	if cfg.LoraBackend.Type != TransportMQTT && cfg.LoraBackend.Type != TransportHTTP {
		typ = TransportMQTT
	} else {
		typ = cfg.LoraBackend.Type
	}

	serv := &Server{}
	if typ == TransportMQTT {
		devs, err := storage.GetDevicesEUI()
		if err != nil {
			return nil, err
		}
		serv.backend = mqtt.NewBackend(cfg.LoraBackend.Mqtt, devs)
	} else {
		var addr string
		if cfg.LoraBackend.HTTPPort > 0 {
			addr = fmt.Sprintf(":%d", cfg.LoraBackend.HTTPPort)
		} else {
			addr = ":8080"
		}
		httpserv := http.New(addr)

		serv.backend = httpserv
	}

	return serv, nil
}

// Start server start
func (s *Server) Start() {
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		s.handleUplinks(&s.wg)
	}()
}

// Stop 关闭相关资源
func (s *Server) Stop() error {
	if err := s.backend.Close(); err != nil {
		return err
	}
	log.Info("waiting for pending actions to complete")
	s.wg.Wait()
	return nil
}

// OnDeviceChange notice backend subscribe/unsubscribe message
func (s *Server) OnDeviceChange(notice map[string]bool) {
	s.backend.Notice(notice)
}

// HandleUplinks 处理lora上行数据
func (s *Server) handleUplinks(wg *sync.WaitGroup) {
	for uplink := range s.backend.RXPacketChan() {
		go func(uplink backend.DataUpPayloadChan) {
			wg.Add(1)
			defer wg.Done()
			if err := handleUplink(uplink); err != nil {
				log.WithFields(log.Fields{
					"device": uplink.DevEUI,
					"data":   hex.EncodeToString(uplink.Data),
				}).Errorf("process device uplink data error: %s", err)
			}
		}(uplink)
	}

}

func handleUplink(data backend.DataUpPayloadChan) error {
	if data.Data[0] == 0xAA {
		ag := &protocol.Angus{}
		err := ag.Unmarshal(data.Data)
		if err != nil {
			return err
		}

		if err := createTrack(data, ag); err != nil {
			return err
		}

		if err := createState(data, ag); err != nil {
			return err
		}
	}
	return nil
}

func createTrack(data backend.DataUpPayloadChan, ag *protocol.Angus) error {
	track := storage.DeviceTrack{
		DeviceEUI: data.DevEUI,
		CreatedAt: time.Now(),
		Location: storage.GPSPoint{
			Latitude:  ag.Latitude,
			Longitude: ag.Longitude,
		},
		Altitude: ag.Altitude,
	}

	err := storage.CreateDeviceTrack(track)
	if err != nil {
		return err
	}

	return nil
}

func createState(data backend.DataUpPayloadChan, ag *protocol.Angus) error {
	now := time.Now()
	ds := storage.DeviceState{
		DeviceEUI:  data.DevEUI,
		LastSeenAt: &now,
		Location: &storage.GPSPoint{
			Latitude:  ag.Latitude,
			Longitude: ag.Longitude,
		},
	}
	var st = state{
		ID: data.DevEUI.String(),
		Point: point{
			ID: data.DevEUI.String(),
			Geometry: geometry{
				Type:        "Point",
				Coordinates: []float64{ag.Longitude, ag.Latitude},
			},
		},
		Prop: map[string]interface{}{
			"设备ID": data.DevEUI.String(),
		},
		Sensor: map[string]interface{}{
			"速度":  ag.Speed,
			"方位角": ag.Azimuth,
			"海拔":  ag.Altitude,
		},
	}
	if ag.DataField != nil {
		switch v := ag.DataField.(type) {
		case *protocol.AngusAlert:
			if v.SOS {
				st.Sensor["警报"] = "SOS"
			}
			if v.LowBattery {
				st.Sensor["警报"] = "低电压"
			}
			if v.Remove {
				st.Sensor["警报"] = "设备摘除"
			}
		case *protocol.AngusSensor:
			st.Sensor["步数"] = v.StepNumber
			st.Sensor["业务ID"] = v.BusinessID
			st.Sensor["电量百分比"] = fmt.Sprintf("%d%%", v.Power)
		case *protocol.AngusHeartbeat:
			st.Sensor["步数"] = v.StepNumber
			st.Sensor["业务ID"] = v.BusinessID
		}
	}
	b, _ := json.Marshal(st)
	ds.Detail = b

	err := storage.CreateAndUpdateState(ds)
	return err
}

type state struct {
	ID     string                 `json:"id"`
	Point  point                  `json:"point"`
	Prop   map[string]interface{} `json:"prop"`
	Sensor map[string]interface{} `json:"sensor"`
}

// Point vbase point
type point struct {
	ID       string   `json:"id"`
	Geometry geometry `json:"geometry"`
}

// Geometry vbase geometry
type geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
