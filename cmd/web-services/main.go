package webservices

import (
	"context"
	"log"

	"segokuning/api/handlers"
	"segokuning/api/responses"
	"segokuning/api/routes"
	"segokuning/configs"
	"segokuning/db/connections"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Run() {
	app := fiber.New(
		fiber.Config{
			StrictRouting:     true,
			EnablePrintRoutes: true,
			CaseSensitive:     true,
		},
	)

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	dbPool, err := connections.NewPgConn(config)
	if err != nil {
		log.Fatalf("failed open connection to db: %v", err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("FAILED PING TO DB: %v", err)
	}

	deps := handlers.Dependencies{
		Cfg:    config,
		DbPool: dbPool,
	}

	// load Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// register route in another package
	routes.RouteRegister(app, deps)

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return responses.ReturnTheResponse(c, true, int(404), "Not Found", nil)
	})

	// Here we go!
	log.Fatalln(app.Listen(":" + config.APPPort))
}
