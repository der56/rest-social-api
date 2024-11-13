package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error, statusCode int) {
	c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
}

// Specific errors
var (
	ErrEmailTaken            = errors.New("Email is already taken.")
	ErrPasswordHash          = errors.New("Error generating password hash.")
	ErrSaveUser              = errors.New("Error saving user.")
	ErrDecodeRequest         = errors.New("Error decoding request body.")
	ErrBadRequest            = errors.New("Bad request.")
	ErrInternalServerError   = errors.New("Internal server error.")
	ErrUsernameTaken         = errors.New("Username is already taken.")
	ErrEmailAndUsernameTaken = errors.New("Email and username are already taken.")
	ErrGenerateJWT           = errors.New("Error generating JWT.")
	ErrTokenNotRenewable     = errors.New("Error renewing token.")
	ErrInvalidPassword       = errors.New("Error invalid password.")
	ErrFindUser              = errors.New("Error finding user.")
	ErrUserNotFound          = errors.New("Error user not found.")
	ErrMissingFields         = errors.New("Error missing fields.")
	ErrUpdatePassword        = errors.New("Error updating password")
	ErrUnauthorized          = errors.New("Error you are not allowed to do that.")
	ErrUpdateUsername        = errors.New("Error updating username")
	ErrInvalidPasswordType   = errors.New("Error. Invalid password type")
	ErrInvalidEmailType      = errors.New("Error. Invalid email type")
	ErrInvalidLoginParams    = errors.New("Error invalid login params")
)
