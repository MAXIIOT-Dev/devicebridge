package lorahandler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/maxiiot/vbaseBridge/storage"

	"github.com/maxiiot/vbaseBridge/backend/mqtt"

	log "github.com/sirupsen/logrus"
)

// Server 后台服务
type Server struct {
	backend *mqtt.Backend
	wg      sync.WaitGroup
}

// NewServer return server point
func NewServer(backend *mqtt.Backend) *Server {
	return &Server{backend: backend}
}

// Start server start
func (s *Server) Start() {
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		s.HandleUplinks(&s.wg)
	}()
}

// Stop 关闭相关资源
func (s *Server) Stop() error {
	if err := s.backend.Close(); err != nil {
		return fmt.Errorf("close backend mqtt error: %s", err)
	}
	log.Info("waiting for pending actions to complete")
	s.wg.Wait()
	return nil
}

// HandleUplinks 处理lora上行数据
func (s *Server) HandleUplinks(wg *sync.WaitGroup) {
	for uplink := range s.backend.RXPacketChan() {
		go func(uplink mqtt.DataUpPayloadChan) {
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

func handleUplink(data mqtt.DataUpPayloadChan) error {
	if data.Data[0] == 0xAA {
		ag := &Angus{}
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

func createTrack(data mqtt.DataUpPayloadChan, ag *Angus) error {
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

func createState(data mqtt.DataUpPayloadChan, ag *Angus) error {
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
	}
	if ag.DataField != nil {
		detail := make(map[string]interface{})
		switch v := ag.DataField.(type) {
		case *AngusAlert:
			if v.SOS {
				detail["警报"] = "SOS"
			}
			if v.LowBattery {
				detail["警报"] = "低电压"
			}
			if v.Remove {
				detail["警报"] = "设备摘除"
			}
		case *AngusSensor:
			detail["步数"] = v.StepNumber
			detail["业务ID"] = v.BusinessID
			detail["电量百分比"] = fmt.Sprintf("%d%%", v.Power)
		case *AngusHeartbeat:
			detail["步数"] = v.StepNumber
			detail["业务ID"] = v.BusinessID
		}
		st.Sensor = detail
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
