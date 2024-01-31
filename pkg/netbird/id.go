package netbird

import (
	"crypto/sha256"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func GetRandomID(length int) string {
	if length > 64 {
		length = 64
	}
	uuid := uuid.NewV4()
	h := sha256.New()
	h.Write(uuid.Bytes())
	bs := h.Sum(nil)
	shortSha := fmt.Sprintf("%x", bs)[:length]
	return shortSha
}
