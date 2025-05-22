package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

const defaultShortCodeLength = 7

func GenerateShortCode(originalURL string) string {
	if originalURL == "" {
		return ""
	}

	hasher := sha256.New()
	hasher.Write([]byte(originalURL))
	hashBytes := hasher.Sum(nil)

	hashString := hex.EncodeToString(hashBytes)

	if len(hashString) < defaultShortCodeLength {
		return hashString
	}
	return hashString[:defaultShortCodeLength]
}
