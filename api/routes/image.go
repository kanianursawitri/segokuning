package routes

import (
	"segokuning/api/handlers"
	"segokuning/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ImageRoutes(app *fiber.App, h handlers.ImageUploader) {
	app.Post("/v1/image", middleware.JWTAuth(), h.Upload)
}
