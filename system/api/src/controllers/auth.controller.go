package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"social_api/libs"
	"social_api/models"
	"social_api/schemas"
	"social_api/utils"

	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate
var translator ut.Translator

func RegisterHandler(c *fiber.Ctx) error {
	var requestBody schemas.RegistrationRequest
	if err := json.Unmarshal([]byte(c.Body()), &requestBody); err != nil {
		utils.HandleError(c, utils.ErrDecodeRequest, http.StatusBadRequest)
		return errors.New("Error decoding request body.")
	}

	username, email, password, firstname, lastname := requestBody.Username, requestBody.Email, requestBody.Password, requestBody.Firstname, requestBody.Lastname

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(c, utils.ErrPasswordHash, http.StatusInternalServerError)
		fmt.Println("Error generating password hash:", err)
		return errors.New("Error generating password hash.")
	}

	newUser := models.User{
		Username:  username,
		Email:     email,
		Password:  string(passwordHash),
		FirstName: firstname,
		LastName:  lastname,
	}

	newUserID, errSavingUser := utils.SaveUser(newUser)

	if errSavingUser != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error saving user": errSavingUser.Error(),
		})
	}

	token, err := libs.GenerateJWT(newUserID.String)
	if err != nil {
		// Manejar el error
		utils.HandleError(c, utils.ErrGenerateJWT, http.StatusInternalServerError)
		c.Status(http.StatusInternalServerError).Send([]byte("Error generating JWT."))
		return errors.New("Error generando token JWT")
	}

	response := map[string]interface{}{
		"message": "Registration successful",
		"user":    utils.UserWithoutPassword(newUser, newUserID.String),
	}

	c.Set("session", token)
	return c.Status(200).JSON(response)
}

func LoginHandler(c *fiber.Ctx) error {
	var requestBody schemas.LoginRequest

	if err := json.Unmarshal([]byte(c.Body()), &requestBody); err != nil {
		utils.HandleError(c, utils.ErrDecodeRequest, http.StatusBadRequest)
		return errors.New("Error decoding request body.")
	}

	if requestBody.Username == "" && requestBody.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Either 'Username' or 'Email' is required.",
		})
	}

	expectedFields := map[string]bool{
		"Username": true,
		"Email":    true,
		"Password": true,
	}

	unknownFields := make(map[string]interface{})
	if err := json.Unmarshal([]byte(c.Body()), &unknownFields); err != nil {
		utils.HandleError(c, utils.ErrDecodeRequest, http.StatusBadRequest)
		return errors.New("Error decoding request body.")
	}

	// Identificar los campos desconocidos
	unknownFieldList := make([]string, 0)
	for field := range unknownFields {
		if !expectedFields[field] {
			unknownFieldList = append(unknownFieldList, field)
		}
	}

	if len(unknownFieldList) > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":          "Unknown fields in the request.",
			"unknown_fields": unknownFieldList,
		})
	}

	var user *models.User
	var err error

	if requestBody.Username != "" {
		user, err = utils.FindUserByEmailOrUsername("", requestBody.Username)
	} else {
		user, err = utils.FindUserByEmailOrUsername(requestBody.Email, "")
	}

	if err != nil {
		utils.HandleError(c, utils.ErrFindUser, http.StatusInternalServerError)
		return errors.New("Error searching the user.")
	}

	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found.",
		})
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Wrong password",
			})
		}
	}
	var userID string
	if user.ID.Valid {
		userID = user.ID.String
	} else {
		userID = ""
	}

	token, err := libs.GenerateJWT(userID)
	if err != nil {
		utils.HandleError(c, utils.ErrGenerateJWT, http.StatusInternalServerError)
		c.Status(http.StatusInternalServerError).Send([]byte("Error generating JWT token"))
		return errors.New("Error generating JWT token.")
	}

	response := map[string]interface{}{
		"message": "Login successful.",
		"user":    utils.UserWithoutPassword(*user, userID),
	}

	c.Set("session", token)
	return c.Status(200).JSON(response)
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("session")
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Logout successful.",
	})
}
