package routes

import (
	"shopifyx/api/handlers"
	"shopifyx/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userHandler handlers.User) {
	g := app.Group("/v1/user")
	g.Post("/register", userHandler.Register)
	g.Post("/login", userHandler.Login)
	// protected routes
	g.Patch("", middleware.JWTAuth(), userHandler.UpdateAccount)
	g.Post("/link/email", middleware.JWTAuth(), userHandler.UpdateEmail)
	g.Post("/link/phone", middleware.JWTAuth(), userHandler.UpdatePhone)
}
