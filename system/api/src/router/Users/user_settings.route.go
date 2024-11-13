package router

import (
	"social_api/controllers"
	"social_api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupUserSettiingsRoutes(app *fiber.App) {
	userSettingsRouter := app.Group("/api")

	userSettingsRouter.Post("/update-username", middlewares.RenewJWTMiddleware, controllers.UpdateUsernameHandler)
	userSettingsRouter.Post("/update-password", middlewares.RenewJWTMiddleware, controllers.UpdatePasswordHandler)
}
