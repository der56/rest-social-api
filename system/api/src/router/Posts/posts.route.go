package router

import "github.com/gofiber/fiber/v2"

func SetUpPostsRoutes(app *fiber.App) {
	postsRouter := app.Group("/api")

	postsRouter.Get("/posts", func(c *fiber.Ctx) error {
		return nil
	})

	// On development
}
