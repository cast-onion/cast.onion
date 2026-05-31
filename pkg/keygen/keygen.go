package keygen

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const (
	stationKeyLen  = 32
	accessTokenLen = 48
)

func GenerateStationKey() (raw string, hash string, err error) {
	return generate("sk_", stationKeyLen)
}

func GenerateAccessToken() (raw string, hash string, err error) {
	return generate("at_", accessTokenLen)
}

func generate(prefix string, n int) (raw, hash string, err error) {
	b := make([]byte, n)
	if _, err = rand.Read(b); err != nil {
		return
	}
	raw = prefix + base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(sum[:])
	return
}

func HashRaw(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func GenerateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
