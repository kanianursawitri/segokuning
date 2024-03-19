package routes

import (
	"shopifyx/api/handlers"
	"shopifyx/db/functions"

	"github.com/gofiber/fiber/v2"
)

func RouteRegister(app *fiber.App, deps handlers.Dependencies) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	userHandler := handlers.User{
		Database: functions.NewUser(deps.DbPool, deps.Cfg),
	}

	postHandler := handlers.Post{
		Database: functions.NewPost(deps.DbPool, deps.Cfg),
	}

	commentHandler := handlers.Comment{
		Database:       functions.NewComment(deps.DbPool, deps.Cfg),
		PostDatabase:   functions.NewPost(deps.DbPool, deps.Cfg),
		FriendDatabase: functions.NewFriend(deps.DbPool, deps.Cfg),
	}

	UserRoutes(app, userHandler)
	PostRoutes(app, postHandler)
	CommentRoutes(app, commentHandler)
}
