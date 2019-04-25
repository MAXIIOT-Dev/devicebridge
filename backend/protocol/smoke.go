package protocol

import (
	"fmt"
)

// SmokeAlarm 烟雾报警类型
type SmokeAlarm [2]byte

var (
	// SmokeAlarmSmoke 烟雾报警
	SmokeAlarmSmoke SmokeAlarm = [2]byte{0x00, 0x01}
	// SmokeAlarmHightTemp 高温报警
	SmokeAlarmHightTemp = [2]byte{0x00, 0x02}
	// SmokeAlarmSmokeAndHightTemp 烟雾和高温报警
	SmokeAlarmSmokeAndHightTemp = [2]byte{0x00, 0x04}
	// SmokeAlarmSensorFail 烟雾传感器故障
	SmokeAlarmSensorFail = [2]byte{0x00, 0x08}
	// SmokeAlarmHightTempSensorFail 高温传感器故障
	SmokeAlarmHightTempSensorFail = [2]byte{0x00, 0x10}
	// SmokeAlarmHSSensorFail 高温和烟雾传感器故障
	SmokeAlarmHSSensorFail = [2]byte{0x00, 0x20}
	// SmokeAlarmLowerEle 系统低电量故障
	SmokeAlarmLowerEle = [2]byte{0x00, 0x40}
	// SmokeAlarmLowSensitivity 烟雾传感器灵敏度过低故障
	SmokeAlarmLowSensitivity = [2]byte{0x00, 0x80}
	// SmokeAlarmHighSensitivity 烟雾传感器灵敏度过高故障
	SmokeAlarmHighSensitivity = [2]byte{0x01, 0x00}
)

func (sa SmokeAlarm) String() string {
	switch sa {
	case SmokeAlarmSmoke:
		return "Smoke alarm"
	case SmokeAlarmHightTemp:
		return "High temperature alarm"
	case SmokeAlarmSmokeAndHightTemp:
		return "Smoke and high temperature alarm"
	case SmokeAlarmSensorFail:
		return "Smoke sensor failure"
	case SmokeAlarmHightTempSensorFail:
		return "High temperature sensor failure"
	case SmokeAlarmHSSensorFail:
		return "High temperature and smoke sensor failure"
	case SmokeAlarmLowerEle:
		return "System low battery failure"
	case SmokeAlarmLowSensitivity:
		return "Smoke sensor sensitivity is too low"
	case SmokeAlarmHighSensitivity:
		return "Smoke sensor sensitivity is too high"
	default:
		return "unknown alert"
	}
}

func (sa SmokeAlarm) MarshalText() ([]byte, error) {
	return []byte(sa.String()), nil
}

// Smoke smoke device handler
type Smoke struct {
	IsHeartBeat bool        `json:"is_heartbeat"`
	Alarm       *SmokeAlarm `json:"alarm,omitempty"`
}

// Unmarshal decode data to struct
// demo: 1800000601020200406381
func (s *Smoke) Unmarshal(data []byte) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()

	start := 0
	FunctionCodeNum := data[start]
	start++
	for i := 0; i < int(FunctionCodeNum); i++ {
		functionCodeID := data[start]
		switch functionCodeID {
		case 0x00: // 心跳包
			s.IsHeartBeat = true
			start++
			dataLen := int(data[start])
			start++
			start += dataLen

		case 0x01: // 下行应答
			start++
			dataLen := int(data[start])
			start++
			start += dataLen
		case 0x02: // 烟雾报警上报
			start++
			dataLen := int(data[start])
			start++
			data := data[start : start+dataLen]
			alarm := SmokeAlarm{}
			copy(alarm[:], data)
			s.Alarm = &alarm
			start += dataLen
		default:
			start++
			dataLen := int(data[start])
			start++
			start += dataLen
		}
	}
	return
}
