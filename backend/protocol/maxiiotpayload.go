package protocol

import (
	"errors"
	"fmt"
)

// Payload define data field interface
type Payload interface {
	Unmarshal([]byte) error
}

// TransCode 传输码
type TransCode struct {
	Direction string `json:"direction"` // 方向
	Major     byte   `json:"major"`     // 主版本号
	Minor     byte   `json:"minor"`     // 次版本号
}

// Unmarshal TransCode unamrshal
func (tc *TransCode) Unmarshal(b byte) error {
	if b&1<<7 == 1 {
		tc.Direction = "Server2Dev"
	} else {
		tc.Direction = "Dev2Server"
	}
	tc.Major = b & (1<<6 | 1<<5 | 1<<4)
	tc.Minor = b & (1<<3 | 1<<2 | 1)
	return nil
}

// MaxiiotPayload define maxiiot device payload
type MaxiiotPayload struct {
	Header     byte       // 帧起始标志
	TransCode  *TransCode // 传输码
	DeviceID   [2]byte    // 设备ID
	SensorData Payload    // 数据域
	CRC        byte       // 校验码
	End        byte       // 帧结束标志
}

// Unmarshal unmarshal maxiiot device
func (p *MaxiiotPayload) Unmarshal(b []byte) (err error) {
	defer func() {
		if res := recover(); res != nil {
			err = fmt.Errorf("panic: %v", res)
		}
	}()

	length := len(b)
	if length < 9 {
		return errors.New("unspoorts maxiiot device protocol")
	}
	flag := 0
	p.Header = b[flag]
	flag++
	tc := &TransCode{}
	tc.Unmarshal(b[flag])
	p.TransCode = tc
	flag++
	copy(p.DeviceID[:], b[flag:flag+2])
	flag += 2
	switch p.DeviceID {
	case [2]byte{0x00, 0x06}:
		smoke := &Smoke{}
		if err := smoke.Unmarshal(b[flag : length-2]); err != nil {
			return err
		}
		p.SensorData = smoke
	default:
		return errors.New("unspoorts maxiiot device protocol")
	}

	return nil
}
