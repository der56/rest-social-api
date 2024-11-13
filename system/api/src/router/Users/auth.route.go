package router

import (
	"social_api/controllers"
	"social_api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthenticationRoutes(app *fiber.App) {
	authRouter := app.Group("/api")

	authRouter.Post("/register", middlewares.ValidateRegisterSchema, controllers.RegisterHandler)
	authRouter.Post("/login", middlewares.ValidateLoginSchema, controllers.LoginHandler)
	authRouter.Post("/logout", controllers.LogoutHandler)
	authRouter.Get("/profile", middlewares.RenewJWTMiddleware, controllers.PrivateProfileHandler)
}
