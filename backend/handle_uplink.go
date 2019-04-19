package backend

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/maxiiot/vbaseBridge/backend/protocol"
	"github.com/maxiiot/vbaseBridge/storage"
)

// HandleUplink handle uplink data
func HandleUplink(data DataUpPayloadChan) error {
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

func createTrack(data DataUpPayloadChan, ag *protocol.Angus) error {
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

func createState(data DataUpPayloadChan, ag *protocol.Angus) error {
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
			"设备ID":   data.DevEUI.String(),
			"最后上传时间": now.Format(time.RFC3339),
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
