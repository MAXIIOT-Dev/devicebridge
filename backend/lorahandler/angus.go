package lorahandler

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Payload define data field interface
type Payload interface {
	Unmarshal([]byte) error
}

// Angus 牛羊定位器
type Angus struct {
	FrameHeader     uint8     `json:"frame_header"` // 帧头
	UTC             time.Time `json:"utc"`          // UTC时间
	OriginLatitude  uint32    `json:"-"`            // 原始纬度
	OriginLongitude uint32    `json:"-"`            // 原始经度
	Latitude        float64   `json:"latitude"`     // 纬度
	Longitude       float64   `json:"longitude"`    // 经度
	Speed           uint8     `json:"speed"`        // 速度
	Azimuth         uint16    `json:"azimuth"`      // 方位角
	Altitude        uint16    `json:"altitude"`     // 海拔
	Code            uint8     `json:"code"`         // 功能码
	DataLen         uint8     `json:"data_len"`     // 数据长度
	DataField       Payload   `json:"data_field"`   // 数据域
	CRC             uint8     `json:"crc"`          // 校验码
}

// Unmarshal 数据解析
func (a *Angus) Unmarshal(data []byte) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()

	length := len(data)
	if length <= 21 {
		return fmt.Errorf("Angus payload length must >=21")
	}

	if data[0] != 0xAA {
		return fmt.Errorf("Angus payload start with 0xAA")
	}
	st := 0
	a.FrameHeader = data[st]
	st += 1

	utc := binary.BigEndian.Uint32(data[st : st+4])
	st += 4
	a.UTC = time.Unix(int64(utc), 0)

	a.OriginLatitude = binary.BigEndian.Uint32(data[st : st+4])
	st += 4
	a.Latitude = float64(a.OriginLatitude) / 1000000.

	a.OriginLongitude = binary.BigEndian.Uint32(data[st : st+4])
	st += 4
	a.Longitude = float64(a.OriginLongitude) / 1000000.

	a.Latitude, a.Longitude = gps84ToGcj02(a.Latitude, a.Longitude)

	a.Speed = data[st]
	st += 1
	a.Azimuth = binary.BigEndian.Uint16(data[st : st+2])
	st += 2
	a.Altitude = binary.BigEndian.Uint16(data[st : st+2])
	st += 2
	a.Code = data[st]
	st += 1
	a.DataLen = data[st]
	st += 1

	switch a.Code {
	case 0x01:
		if length != 22 {
			return fmt.Errorf("报警功能数据长度应为22")
		}
		alert := &AngusAlert{}
		alert.Unmarshal(data[st : st+1])
		a.DataField = alert
	case 0x02:
		if length != 28 {
			return fmt.Errorf("传感器信息数据长度应为28")
		}
		sensor := &AngusSensor{}
		sensor.Unmarshal(data[st : st+7])
		a.DataField = sensor
	case 0x03:
		if length != 27 {
			return fmt.Errorf("心跳包信息数据长度应为27")
		}
		hb := &AngusHeartbeat{}
		hb.Unmarshal(data[st : st+6])
		a.DataField = hb
	}
	a.CRC = data[length-1]
	return nil
}

// AngusAlert 报警提醒
type AngusAlert struct {
	SOS        bool `json:"sos"`         // SOS 警报
	LowBattery bool `json:"low_battery"` // 低电压警报
	Remove     bool `json:"remove"`      // 摘除警报
}

// Unmarshal  AngusAlert unmarshal
func (alert *AngusAlert) Unmarshal(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("报警提醒长度为一个字节")
	}
	if data[0] == 0x01 {
		alert.SOS = true
	}
	if data[0] == 0x02 {
		alert.LowBattery = true
	}
	if data[0] == 0x04 {
		alert.Remove = true
	}
	return nil
}

// AngusSensor 传感器信息
type AngusSensor struct {
	StepNumber uint16 `json:"step_number"` // 步数
	BusinessID uint32 `json:"business_id"` // 业务 ID
	Power      uint8  `json:"power"`       // 电量百分比
}

// Unmarshal AngusSensor Unmarshal
// demo: aa5cac117c0158a42e06ca2e5c01015bffea01010466
func (as *AngusSensor) Unmarshal(data []byte) error {
	if len(data) != 7 {
		return fmt.Errorf("传感器信息长度为7个字节")
	}
	as.StepNumber = binary.BigEndian.Uint16(data[:2])
	as.BusinessID = binary.BigEndian.Uint32(data[2:6])
	as.Power = data[6]
	return nil
}

// AngusHeartbeat 心跳包信息
type AngusHeartbeat struct {
	StepNumber uint16 `json:"step_number"` // 步数
	BusinessID uint32 `json:"business_id"` // 业务 ID
}

// Unmarshal AngusHeartbeat unmarshal
func (ah *AngusHeartbeat) Unmarshal(data []byte) error {
	if len(data) != 6 {
		return fmt.Errorf("心跳信息长度为6个字节")
	}
	ah.StepNumber = binary.BigEndian.Uint16(data[:2])
	ah.BusinessID = binary.BigEndian.Uint32(data[2:6])

	return nil
}
