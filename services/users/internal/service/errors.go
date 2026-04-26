package service

import (
	"errors"
	"fmt"
)

var ErrPasswordTooShort = errors.New("password is too short")
var ErrAuthStateNil = errors.New("internal auth state is nil")

func registerError(username string, err error) error {
	return fmt.Errorf("auth service: register for %s failed: %w", username, err)
}
