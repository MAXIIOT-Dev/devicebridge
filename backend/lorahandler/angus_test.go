package lorahandler

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

func TestAngusUnamrshal(t *testing.T) {
	data, err := hex.DecodeString("aa5cac117c0158a42e06ca2e5c01015bffea01010466")
	if err != nil {
		t.Error("decode string data error:", err)
	}
	ang := Angus{}
	ang.Unmarshal(data)
	b, _ := json.MarshalIndent(ang, "", " ")
	t.Log(string(b))
}
