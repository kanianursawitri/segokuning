package routes

import (
	"segokuning/api/handlers"
	"segokuning/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func FriendRoutes(app *fiber.App, friendHandler handlers.Friend) {
	g := app.Group("/v1/friend")
	g.Get("", middleware.JWTAuth(), friendHandler.GetFriends)
	g.Post("", middleware.JWTAuth(), friendHandler.AddFriend)
	g.Delete("", middleware.JWTAuth(), friendHandler.DeleteFriend)
}
