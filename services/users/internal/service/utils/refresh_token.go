package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		return "", fmt.Errorf("utils: generate refresh token failed: rand.Read: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
