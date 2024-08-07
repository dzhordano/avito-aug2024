package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	UserType UserType
}

type UserType string

const (
	UserTypeClient    UserType = "client"
	UserTypeModerator UserType = "moderator"
)

func (u UserType) Validate() bool {
	return u != "" && (u == UserTypeClient || u == UserTypeModerator)
}

func (u UserType) String() string {
	return string(u)
}
