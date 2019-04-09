package utils

import (
	"testing"
)

func Test_Hash(t *testing.T) {
	hash, err := Hash("admin")
	if err != nil {
		t.Error("hash error:", err)
	}
	t.Log(hash)
}
