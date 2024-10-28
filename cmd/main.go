package main

import (
	"fmt"
	"log"
	"payment-microservice/config"
	"payment-microservice/internal/handlers"

	"payment-microservice/pkg/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.GetConfig()

	postgres.ConnectDb()

	app := fiber.New()

	app.Use(logger.New())

	handlers.SetupRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.Config.Port)))
}
