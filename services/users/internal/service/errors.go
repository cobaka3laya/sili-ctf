package service

import (
	"errors"
	"fmt"
)

var ErrUserStateNil = errors.New("internal user state is nil")
var ErrRefreshTokenStatenil = errors.New("internal refresh token state is nil")
var ErrAuthDataStateNil = errors.New("internal auth state is nil")

var ErrPasswordIncorrect = errors.New("password is incorrect")
var ErrPasswordTooShort = errors.New("password is too short")

func registerError(username string, err error) error {
	return fmt.Errorf("auth service: register for %s failed: %w", username, err)
}

func loginError(username string, err error) error {
	return fmt.Errorf("auth service: login for %s failed: %w", username, err)
}
