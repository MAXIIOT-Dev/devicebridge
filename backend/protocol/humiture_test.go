package protocol

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

func Test_Humitureunmarshal(t *testing.T) {
	data, err := hex.DecodeString("ff0200015cc11b9401011706013b014aff")
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
