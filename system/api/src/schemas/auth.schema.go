package schemas

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var (
	ErrValidation = errors.New("Validation error")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type RegistrationRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=20"`
	Email     string `json:"email" validate:"required,email,max=35"`
	Password  string `json:"password" validate:"required,min=6,max=45"`
	Firstname string `json:"firstname" validate:"required,max=18"`
	Lastname  string `json:"lastname" validate:"required,max=50"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUsernameRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" validate:"required,min=6"`
}

func Validate(request interface{}) error {
	return validate.Struct(request)
}
