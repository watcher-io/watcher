package utility

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
