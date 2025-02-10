package sec

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(val string) string {
	hash := sha256.Sum256([]byte(val))
	return hex.EncodeToString(hash[:])
}
