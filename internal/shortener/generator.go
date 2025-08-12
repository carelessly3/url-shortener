package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

const codeLen = 6

// GenerateCode returns a URL-safe alphanumeric code of length ~codeLen
func GenerateCode() (string, error) {
	// generate raw random bytes and base64-url encode, then strip non-alnum and trim
	b := make([]byte, codeLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := base64.RawURLEncoding.EncodeToString(b) // URL-safe
	// keep only alnum and limit length
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, s)
	if len(s) < codeLen {
		return "", errors.New("generated code too short")
	}
	return s[:codeLen], nil
}
