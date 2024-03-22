package routes

import (
	"segokuning/api/handlers"
	"segokuning/db/functions"
	"segokuning/internal/utils"

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
		PostDatabase:   functions.NewPost(deps.DbPool, deps.Cfg),
		FriendDatabase: functions.NewFriend(deps.DbPool, deps.Cfg),
	}

	imageUploaderHandler := handlers.ImageUploader{
		Uploader: utils.NewImageUploader(deps.Cfg),
	}

	friendHandler := handlers.Friend{
		Database: functions.NewFriend(deps.DbPool, deps.Cfg),
	}

	ImageRoutes(app, imageUploaderHandler)
	UserRoutes(app, userHandler)
	PostRoutes(app, postHandler)
	CommentRoutes(app, commentHandler)
	FriendRoutes(app, friendHandler)
}
