package protocol

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

func Test_Humitureunmarshal(t *testing.T) {
	data, err := hex.DecodeString("ff0200015cc056d001011b0401460151ff")
	if err != nil {
		t.Error("decode data error:", err)
	}
	var hums Humitures
	err = hums.Unmarshal(data)
	if err != nil {
		t.Error(err)
	}
	j, _ := json.MarshalIndent(hums, "", " ")
	t.Log(string(j))
}
