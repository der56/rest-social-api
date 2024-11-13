package router

import (
	posts "social_api/router/Posts"
	social "social_api/router/Social"
	users "social_api/router/Users"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	users.SetupAuthenticationRoutes(app)
	posts.SetUpPostsRoutes(app)
	social.SetUpSocialRoutes(app)
	users.SetupUserSettiingsRoutes(app)
}
