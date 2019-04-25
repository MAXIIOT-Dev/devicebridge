package protocol

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

func TestSmokeUnmarshal(t *testing.T) {
	data, err := hex.DecodeString("1800000601020200406381")
	if err != nil {
		t.Error("decode data error:", err)
	}

	s := Smoke{}
	err = s.Unmarshal(data[4:])

	if err != nil {
		t.Error("smoke unmarshal error:", err)
	}
	js, _ := json.MarshalIndent(s, "", " ")
	t.Log(string(js))
}
