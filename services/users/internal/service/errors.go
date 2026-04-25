package service

import "errors"

var ErrPasswordTooShort = errors.New("password is too short")
var ErrAuthStateNil = errors.New("internal auth state is nil")
