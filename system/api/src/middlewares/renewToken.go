package middlewares

import (
	"os"
	"social_api/libs"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func RenewJWTMiddleware(c *fiber.Ctx) error {
	currentTokenString := c.Get("session")

	if currentTokenString == "" {
		return c.Next()
	}

	currentToken, err := jwt.Parse(currentTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_JWT")), nil
	})

	if err != nil {
		// Check if the error is due to token expiration
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Continue with token renewal logic even if it's expired
			} else {
				// Other errors are treated as invalid tokens
				return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
			}
		}
	}

	if currentToken.Valid {
		return c.Next()
	}

	claims, ok := currentToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Extract user ID from claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token: cant extract ID from user"})
	}

	// Renew the token
	_, err = generateAndSetNewToken(c, userID)
	if err != nil {
		return err
	}

	return c.Next()
}

func generateAndSetNewToken(c *fiber.Ctx, userID string) (string, error) {
	newToken, err := libs.GenerateJWT(userID)
	if err != nil {
		return "", c.Status(500).JSON(fiber.Map{"error": "Error while generating new token: " + err.Error()})
	}

	// Set the new token in the response
	c.Set("session", newToken)
	return newToken, nil
}
