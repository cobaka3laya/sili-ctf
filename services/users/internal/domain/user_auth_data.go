package domain

import (
	"errors"
	"strings"
)

type UserAuthData struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	HashedPassword string `json:"hashed_password"`
}

var ErrEmptyHashedPassword = errors.New("hashed password is empty")
var ErrInvalidUserID = errors.New("user ID is invalid")

func NewUserAuthData(userID int64, hashedPassword string) (*UserAuthData, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	if strings.TrimSpace(hashedPassword) == "" {
		return nil, ErrEmptyHashedPassword
	}

	return &UserAuthData{
		UserID:         userID,
		HashedPassword: hashedPassword,
	}, nil
}
