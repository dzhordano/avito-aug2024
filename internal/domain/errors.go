package domain

import (
	"errors"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadySubscribed = errors.New("user already subscribed")
	ErrFlatNotFound          = errors.New("flat not found")
	ErrFlatAlreadyExists     = errors.New("flat already exist")
	ErrHouseNotFound         = errors.New("house not found")
	ErrHouseAlreadyExists    = errors.New("house already exist")
)
