package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int64, ttl time.Duration, cfg *config.SecretsConfig) (string, error) {
	stringUserID := strconv.FormatInt(userID, 10)

	claims := Claims{
		UserID: stringUserID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Subject:   stringUserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(cfg.JWTGenerationKey))

	if err != nil {
		return "", fmt.Errorf("utils: generate JWT failed for user id %d: token.SignedString: %w", userID, signedString)
	}

	return signedString, nil
}

func ParseJWT(tokenStr string, cfg *config.SecretsConfig) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTGenerationKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("utils: parse JWT failed for token string %s: jwt.ParseWithClaims: %w", tokenStr, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
