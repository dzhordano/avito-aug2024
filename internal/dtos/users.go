package dtos

import (
	"errors"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"net/mail"
)

type UserRegisterInput struct {
	Email    string          `json:"email" binding:"required"`
	Password string          `json:"password" binding:"required"`
	UserType domain.UserType `json:"userType" binding:"required"`
}

type UserLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u UserRegisterInput) Validate() error {
	if !u.UserType.Validate() {
		return errors.New("invalid user type")
	}

	if u.Password == "" {
		return errors.New("invalid password")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}

func (u *UserLoginInput) Validate() error {
	if u.Password == "" {
		return errors.New("invalid password")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}
