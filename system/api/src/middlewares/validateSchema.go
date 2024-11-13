package middlewares

import (
	"encoding/json"
	"net/http"
	"social_api/schemas"
	"social_api/utils"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateRegisterSchema(c *fiber.Ctx) error {
	var requestBody schemas.RegistrationRequest

	if err := json.Unmarshal([]byte(c.Body()), &requestBody); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal number into Go struct field") {
			switch {
			case strings.Contains(err.Error(), "username"):
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "You cannot use the type number in username.",
				})
			case strings.Contains(err.Error(), "firstname"):
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "You cannot use the type number in firstname.",
				})
			case strings.Contains(err.Error(), "lastname"):
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "You cannot use the type number in lastname.",
				})
			case strings.Contains(err.Error(), "password"):
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "You cannot use the type number in password.",
				})
			case strings.Contains(err.Error(), "email"):
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "You cannot use the type number in email.",
				})
			default:
				utils.HandleError(c, utils.ErrDecodeRequest, http.StatusBadRequest)
			}
		} else if strings.Contains(err.Error(), "invalid character") {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		} else {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return err
	}

	if err := schemas.Validate(requestBody); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)

			for _, vErr := range validationErr {
				fieldName := vErr.Field()
				tagName := vErr.Tag()

				switch tagName {
				case "required":
					errorMessages[fieldName] = "This field is required."
				case "min":
					errorMessages[fieldName] = "This field must be at least " + vErr.Param() + " characters."
				case "max":
					errorMessages[fieldName] = "This field must be at most " + vErr.Param() + " characters."
				case "email":
					errorMessages[fieldName] = "This field must be a valid email address."
				default:
					errorMessages[fieldName] = "This field is invalid."
				}
			}

			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"errors": errorMessages,
			})
		}
		return err
	}

	return c.Next()
}

func ValidateLoginSchema(c *fiber.Ctx) error {
	var requestBody schemas.LoginRequest

	if err := json.Unmarshal([]byte(c.Body()), &requestBody); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal number into Go struct field") {
			switch {
			case strings.Contains(err.Error(), "username"):
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "You cannot use the type number in username.",
				})
			case strings.Contains(err.Error(), "password"):
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "You cannot use the type number in password.",
				})
			case strings.Contains(err.Error(), "email"):
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "You cannot use the type number in email.",
				})
			default:
				utils.HandleError(c, utils.ErrDecodeRequest, http.StatusBadRequest)
			}
		} else if strings.Contains(err.Error(), "invalid character") {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		} else {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return err
	}

	if err := schemas.Validate(requestBody); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)

			for _, vErr := range validationErr {
				fieldName := vErr.Field()
				tagName := vErr.Tag()

				switch tagName {
				case "required":
					errorMessages[fieldName] = "This field is required."
				case "email":
					errorMessages[fieldName] = "This field must be a valid email address."
				default:
					errorMessages[fieldName] = "This field is invalid."
				}
			}

			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"errors": errorMessages,
			})
		}
		return err
	}

	return c.Next()
}
