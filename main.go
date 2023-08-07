package main

import (
	"brms/config"
	"brms/pkg/middlewares"
	"brms/services/rules_management/handlers"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	//app.Use(middlewares.UndefinedRoutesMiddleware)
	app.Use(middlewares.ErrorMiddleware())

	// register routes
	handlers.Routes(app)

	listenPort := fmt.Sprintf(":%s", config.GetConfig().Port)

	if err := app.Listen(listenPort); err != nil {
		log.Fatalln("application failed to fired up: ", err)
	}
}
