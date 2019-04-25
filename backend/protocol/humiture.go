package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Humiture defines the humiture item
type Humiture struct {
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Electricity float64   `json:"electricity"`
	DateTime    time.Time `json:"date_time"`
	alarm       *humitureAlarm
}

type humitureAlarm struct {
	humidityLow     bool
	humidityHigh    bool
	temperatureLow  bool
	temperatureHigh bool
	electricityLow  bool
}

// Humitures defines humiture list
type Humitures struct {
	Hums []Humiture
}

// Marshal returns humiture encode json
func (h *Humiture) Marshal() []byte {
	b, err := json.Marshal(h)
	if err != nil {
		return nil
	}
	return b
}

func (ha *humitureAlarm) unmarshal(alarm byte) {
	if alarm&0x01 == 0x01 {
		ha.humidityHigh = true
	}
	if alarm&0x02 == 0x02 {
		ha.temperatureHigh = true
	}
	if alarm&0x04 == 0x04 {
		ha.humidityLow = true
	}
	if alarm&0x08 == 0x08 {
		ha.temperatureLow = true
	}
	if alarm&0x10 == 0x10 {
		ha.electricityLow = true
	}
}

// String returns humiture alarm
func (ha *humitureAlarm) String() string {
	buf := bytes.NewBufferString("")
	if ha.humidityHigh {
		buf.WriteString("湿度过高;")
	}
	if ha.humidityLow {
		buf.WriteString("湿度过低;")
	}
	if ha.temperatureHigh {
		buf.WriteString("温度过高;")
	}
	if ha.temperatureLow {
		buf.WriteString("温度过低;")
	}
	if ha.electricityLow {
		buf.WriteString("电量过低;")
	}
	return string(buf.Bytes())
}

// Unmarshal defines unmarshal data to humitures
// ff0200015cc0528d01011b0801460151ff
func (h *Humitures) Unmarshal(b []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unmarshal error: %s", r)
		}
	}()

	var (
		start   int
		lengthB = len(b)
	)

	if lengthB < 10 {
		return errors.Errorf("数据帧长度小于10")
	}
	if b[0] == 0xff && b[1] == 0x02 {
		var (
			tempInt []int8
			tempDec []int8
			hums    []int8
			eles    []int8
		)
		start += 2
		len := int(binary.BigEndian.Uint16(b[start : start+2]))
		start += 2
		h.Hums = make([]Humiture, 0, len)
		tempInt = make([]int8, 0, len)
		tempDec = make([]int8, 0, len)
		hums = make([]int8, 0, len)
		eles = make([]int8, 0, len)

		ti := int64(binary.BigEndian.Uint32(b[start : start+4]))
		start += 4
		templen1 := int(b[start]) + 10
		start++
		templen2 := int(b[start]) + templen1
		start++
		if lengthB < templen2 {
			return errors.Errorf("数据总长度小于温度长度字节")
		}
		var _tempInt int8
		for start < templen1 {
			if b[start]&0xa0 == 0xa0 {
				for i := 0; i < int(b[start]&0x0f); i++ {
					tempInt = append(tempInt, _tempInt)
				}
			} else {
				_tempInt = int8(b[start])
				tempInt = append(tempInt, _tempInt)
			}
			start++
		}

		var _tempDec int8
		for start < templen2 {
			if b[start]&0xa0 == 0xa0 {
				for i := 0; i < int(b[start]&0x0f); i++ {
					tempDec = append(tempDec, _tempDec)
				}
			} else {
				_tempDec = int8(b[start])
				tempDec = append(tempDec, _tempDec)
			}
			start++
		}

		humlen := templen2 + int(b[start]) + 1
		start++
		if lengthB < humlen {
			return errors.Errorf("数据总长度小于湿度长度字节")
		}
		var _hum int8
		for start < humlen {
			if b[start]&0xa0 == 0xa0 {
				for i := 0; i < int(b[start]&0x0f); i++ {
					hums = append(hums, _hum)
				}
			} else {
				_hum = int8(b[start])
				hums = append(hums, _hum)
			}
			start++
		}

		elelen := humlen + int(b[start]) + 1
		var _ele int8
		for start < elelen {
			if b[start]&0xa0 == 0xa0 {
				for i := 0; i < int(b[start]&0x0f); i++ {
					eles = append(eles, _ele)
				}
			} else {
				_ele = int8(b[start])
				eles = append(eles, _ele)
			}
			start++
		}

		for i := 0; i < len; i++ {
			h.Hums = append(h.Hums, Humiture{
				Temperature: float64(tempInt[i]) + float64(tempDec[i])/10.,
				Humidity:    float64(hums[i]),
				Electricity: float64(eles[i]),
				DateTime:    time.Unix(ti, 0),
			})
			ti += 60
		}
	} else if b[0] == 0xff && b[1] == 0x01 {
		start += 2
		var (
			ti      int64
			tempInt int8
			tempDec int8
			hum     int8
			ele     int8
		)
		ti = int64(binary.BigEndian.Uint32(b[start : start+4]))
		start += 4
		tempInt = int8(b[start])
		start++
		tempDec = int8(b[start])
		start++
		hum = int8(b[start])
		start++
		ele = int8(b[start])
		start++
		alarm := b[start]
		alarminfo := &humitureAlarm{}
		alarminfo.unmarshal(alarm)
		h.Hums = []Humiture{
			Humiture{
				Temperature: float64(tempInt) + float64(tempDec)/10.,
				Humidity:    float64(hum),
				Electricity: float64(ele),
				DateTime:    time.Unix(ti, 0),
				alarm:       alarminfo,
			},
		}
	}
	return nil
}
