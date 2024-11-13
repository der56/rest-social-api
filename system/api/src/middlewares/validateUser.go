package middlewares

import (
	"context"
	"net/http"
	"social_api/db"
	"social_api/utils"

	"github.com/gofiber/fiber/v2"
)

func ValidateUser(c *fiber.Ctx) error {
	token := c.Get("session")
	id := c.Params("id")
	if token == "" {
		utils.HandleError(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": "Missing authorization header.",
		})
	}

	userID, err := utils.ExtractUserIDFromToken(token)
	if err != nil {
		utils.HandleError(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return c.Status(401).JSON(err.Error())
	}

	pool := db.Pool
	query := `SELECT username FROM user_profile WHERE id = $1`
	_, err = pool.Exec(context.Background(), query, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found."})
	}

	if userID != id {
		return c.Status(401).JSON(fiber.Map{"error": "You are not allowed to do that."})
	}

	return c.Next()
}
