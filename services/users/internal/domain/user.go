package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cobaka3laya/sili-ctf/services/users/config"
)

type User struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	CreatedAt         time.Time `json:"time"`
	ProfilePictureURL string    `json:"profile_picture_url"`
}

var ErrUsernameTooShort = errors.New("username is too short")
var ErrUsernameTooLong = errors.New("username is too long")

func NewUser(username string, createdAt time.Time, cfg *config.UserConfig) (*User, error) {
	usrName := strings.TrimSpace(username)

	if cfg.Username.MinLength > utf8.RuneCountInString(usrName) {
		return nil, ErrUsernameTooShort
	}

	if cfg.Username.MaxLength < utf8.RuneCountInString(usrName) {
		return nil, ErrUsernameTooLong
	}

	return &User{
		Username:  username,
		CreatedAt: createdAt,
	}, nil
}
