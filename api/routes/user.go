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
	gp := g.Use(middleware.OptionalJWTAuth())
	gp.Post("/link/email", userHandler.UpdateEmail)
	gp.Post("/link/phone", userHandler.UpdatePhone)
}
