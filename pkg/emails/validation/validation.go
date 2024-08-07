package validation

import (
	"errors"
	"net/mail"
)

type EmailValidator interface {
	Validate(email string) error
}

type emailValidator struct{}

func NewEmailValidator() EmailValidator {
	return &emailValidator{}
}

func (v *emailValidator) Validate(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}
