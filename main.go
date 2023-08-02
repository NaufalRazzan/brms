package main

import (
	"brms/config"
	"brms/services/rules_management/handlers"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())

	handlers.Routes(app)

	listenPort := fmt.Sprintf(":%s", config.GetConfig().Port)

	if err := app.Listen(listenPort); err != nil {
		panic(err)
	}

	fmt.Println("Application start and running!")
}
