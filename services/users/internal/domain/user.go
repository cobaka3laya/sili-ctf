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
	ProfilePictureUrl string    `json:"profile_picture_url"`
}

var ErrUsernameTooShort = errors.New("username is too short")
var ErrUsernameTooLong = errors.New("username is too long")

func NewUser(username string, createdAt time.Time) (*User, error) {
	usrName := strings.TrimSpace(username)

	config := config.MustLoad()

	if config.Data.User.Username.MinLength > utf8.RuneCountInString(usrName) {
		return nil, ErrUsernameTooShort
	}

	if config.Data.User.Username.MaxLength < utf8.RuneCountInString(usrName) {
		return nil, ErrUsernameTooLong
	}

	return &User{
		Username:  username,
		CreatedAt: createdAt,
	}, nil
}
