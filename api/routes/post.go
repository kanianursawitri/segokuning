package routes

import (
	"shopifyx/api/handlers"
	"shopifyx/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func PostRoutes(app *fiber.App, postHandler handlers.Post) {
	g := app.Group("/v1/post")
	g.Post("", middleware.JWTAuth(), postHandler.AddPost)
	g.Get("", middleware.JWTAuth(), postHandler.GetPosts)
}
