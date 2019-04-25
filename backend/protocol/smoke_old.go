package protocol

/*
import (
	"encoding/json"
	"errors"
	"fmt"
)

// Smoke smoke device handler
type Smoke struct {
	FuncCodeNum byte // 功能码数量
	HeartBeat   bool
	Alert       *SmokeAlert
}

type SmokeAlert struct {
	Message string
}

// Unmarshal smoke alert unmarshal
func (sa *SmokeAlert) Unamrshal(b []byte) error {
	if len(b) != 2 {
		return errors.New("数据长度为2")
	}
	var data [2]byte
	copy(data[:], b)
	switch data {
	case [2]byte{0x00, 0x01}:
		sa.Message = "Smoke alarm"
	case [2]byte{0x00, 0x02}:
		sa.Message = "High temperature alarm"
	case [2]byte{0x00, 0x04}:
		sa.Message = "Smoke and high temperature alarm"
	case [2]byte{0x00, 0x08}:
		sa.Message = "Smoke sensor failure"
	case [2]byte{0x00, 0x10}:
		sa.Message = "High temperature sensor failure"
	case [2]byte{0x00, 0x20}:
		sa.Message = "High temperature and smoke sensor failure"
	case [2]byte{0x00, 0x40}:
		sa.Message = "System low battery failure"
	case [2]byte{0x00, 0x80}:
		sa.Message = "Smoke sensor sensitivity is too low"
	case [2]byte{0x01, 0x00}:
		sa.Message = "Smoke sensor sensitivity is too high"
	default:
		sa.Message = "unknown alert"
	}
	return nil
}

func (s *Smoke) marshal() []byte {
	b, _ := json.Marshal(s)
	return b
}

// Unmarshal decode data to struct
func (s *Smoke) Unmarshal(b []byte) (err error) {
	defer func() {
		if res := recover(); res != nil {
			err = fmt.Errorf("panic: %v", res)
		}
	}()
	var (
		flag       int
		FuncCodeID byte
		DataLen    byte
	)
	s.FuncCodeNum = b[flag]
	flag++
	for i := 0; i < int(s.FuncCodeNum); i++ {
		FuncCodeID = b[flag]
		flag++
		switch FuncCodeID {
		case 0x00:
			DataLen = b[flag]
			flag += int(DataLen + 1)
			s.HeartBeat = true
		case 0x02:
			DataLen = b[flag]
			if DataLen != 2 {
				return errors.New("烟雾报警数据域长度应为2")
			}
			alert := &SmokeAlert{}
			alert.Unamrshal(b[flag : flag+2])
			flag += int(DataLen + 1)

		}
	}
	return nil
}

*/
