package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"social_api/schemas"
	"social_api/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func PrivateProfileHandler(c *fiber.Ctx) error {
	token := c.Get("session")
	if token == "" {
		utils.HandleError(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": "Missing authorization header.",
		})
	}

	id, err := utils.ExtractUserIDFromToken(token)
	if err != nil {
		utils.HandleError(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := utils.FindUserById(id)
	if err != nil {
		utils.HandleError(c, utils.ErrFindUser, http.StatusBadRequest)
		return errors.New("Error finding the user")
	}

	response := map[string]interface{}{
		"message": user.Username + " Profile",
		"user":    utils.UserWithoutPassword(*user, user.ID.String),
	}

	return c.Status(200).JSON(response)
}

func UpdateUsernameHandler(c *fiber.Ctx) error {
	tokenString := c.Get("session")
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
	}

	id, err := utils.ParseToken(tokenString)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
	var requestBody schemas.UpdateUsernameRequest
	if err := json.Unmarshal([]byte(c.Body()), &requestBody); err != nil {
		utils.HandleError(c, utils.ErrDecodeRequest, http.StatusBadRequest)
		return errors.New("Error decoding request body.")
	}

	type newUsernameType string
	newUsername := requestBody.Username
	if newUsername == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "New username cannot be empty."})
	}

	if err := utils.UpdateUsername(id, newUsername); err != nil {
		utils.HandleError(c, utils.ErrUpdateUsername, http.StatusInternalServerError)
		return errors.New("Error updating username")
	}

	response := map[string]interface{}{
		"message":     "Username updated successfully",
		"newUsername": newUsername,
	}

	return c.Status(http.StatusOK).JSON(response)
}

func UpdatePasswordHandler(c *fiber.Ctx) error {
	tokenString := c.Get("session")
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
	}

	id, err := utils.ParseToken(tokenString)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	var requestBody schemas.UpdatePasswordRequest
	if err := json.Unmarshal([]byte(c.Body()), &requestBody); err != nil {
		utils.HandleError(c, utils.ErrDecodeRequest, http.StatusBadRequest)
		return errors.New("Error decoding request body.")
	}

	password := requestBody.Password
	if password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "New password cannot be empty"})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(c, utils.ErrPasswordHash, http.StatusInternalServerError)
		return errors.New("Error generating password hash.")
	}

	// Actualizar la password del usuario en la base de datos
	if err := utils.UpdatePassword(id, string(passwordHash)); err != nil {
		utils.HandleError(c, utils.ErrUpdateUsername, http.StatusInternalServerError)
		return errors.New("Error updating username")
	}

	// Respuesta exitosa
	response := map[string]interface{}{
		"message": "Password updated successfully.",
	}

	return c.Status(http.StatusOK).JSON(response)
}
