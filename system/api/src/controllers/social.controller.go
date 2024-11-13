package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"social_api/utils"

	"github.com/gofiber/fiber/v2"
)

func ProfileHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No ID provided",
		})
	}
	user, err := utils.FindUserById(id)
	if err != nil {
		utils.HandleError(c, utils.ErrFindUser, http.StatusBadRequest)
		return errors.New("Error finding the user")
	}

	response := map[string]interface{}{
		"message": user.Username + " Profile",
		"user":    utils.UserWithoutPasswordAndEmail(*user, user.ID.String),
	}

	return c.Status(200).JSON(response)
}

func FollowUserHandler(c *fiber.Ctx) error {
	followedID := c.Params("id")

	if followedID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No ID provided",
		})
	}

	token := c.Get("session")
	if token == "" {
		utils.HandleError(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": "Missing authorization header.",
		})
	}

	followerID, err := utils.ExtractUserIDFromToken(token)
	if err != nil {
		utils.HandleError(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	isAlreadyFollowing, err := utils.IsFollowing(followerID, followedID)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	if isAlreadyFollowing {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("User %s is already being followed by user %s", followedID, followerID),
		})
	}

	errorFollowing := utils.FollowUser(followerID, followedID)

	if errorFollowing != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errorFollowing,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": fmt.Sprintf("User %s followed successfully", followedID),
	})
}

func UnFollowUserHandler(c *fiber.Ctx) error {
	followedID := c.Params("id")

	if followedID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No ID provided",
		})
	}

	token := c.Get("session")
	if token == "" {
		utils.HandleError(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": "Missing authorization header.",
		})
	}

	followerID, err := utils.ExtractUserIDFromToken(token)
	if err != nil {
		utils.HandleError(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	isAlreadyNotFollowing, err := utils.IsNotFollowing(followerID, followedID)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	if isAlreadyNotFollowing {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("User %s is already not being followed by user %s", followedID, followerID),
		})
	}

	errorunFollowing := utils.UnFollowUser(followerID, followedID)

	if errorunFollowing != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errorunFollowing,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": fmt.Sprintf("User %s has successfully unfollowed", followedID),
	})
}

func GetFollowersHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	followers, err := utils.GetFollowers(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Coudn't get followers",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"followers": followers,
	})
}
