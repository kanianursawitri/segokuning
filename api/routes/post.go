package routes

import (
	"segokuning/api/handlers"
	"segokuning/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func PostRoutes(app *fiber.App, postHandler handlers.Post) {
	g := app.Group("/v1/post")
	g.Post("", middleware.JWTAuth(), postHandler.AddPost)
	g.Get("", middleware.JWTAuth(), postHandler.GetPosts)
}
