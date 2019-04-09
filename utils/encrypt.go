package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodeMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	str_md5 := hex.EncodeToString(cipherStr)
	return str_md5
}
