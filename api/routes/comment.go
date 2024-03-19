package routes

import (
	"shopifyx/api/handlers"
	"shopifyx/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func CommentRoutes(app *fiber.App, commentHandler handlers.Comment) {
	g := app.Group("/v1/comment")
	g.Post("", middleware.JWTAuth(), commentHandler.AddComment)
}
