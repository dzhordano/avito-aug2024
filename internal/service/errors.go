package service

import "errors"

var (
	ErrFlatAlreadyExists     = errors.New("flat already exists")
	ErrAlreadyModerating     = errors.New("flat is already on moderation")
	ErrHouseAlreadyExists    = errors.New("house already exists")
	ErrUserAlreadySubscribed = errors.New("user already subscribed")
	ErrHouseNotFound         = errors.New("house not found")
	ErrFlatNotFound          = errors.New("flat not found")
)
