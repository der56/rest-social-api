package router

import (
	"social_api/controllers"
	"social_api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetUpSocialRoutes(app *fiber.App) {
	socialsRouter := app.Group("/api")

	socialsRouter.Get("/profile/:id", middlewares.RenewJWTMiddleware, controllers.ProfileHandler)
	socialsRouter.Post("/followuser/:id", middlewares.RenewJWTMiddleware, controllers.FollowUserHandler)
	socialsRouter.Post("/unfollowuser/:id", middlewares.RenewJWTMiddleware, controllers.UnFollowUserHandler)
	socialsRouter.Get("/getfollowers/:id", middlewares.RenewJWTMiddleware, controllers.GetFollowersHandler)
}
