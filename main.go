package main

import (
	"github.com/gofiber/fiber/v2"

	config "go-htmlcsstoimage/configs"
	"go-htmlcsstoimage/routes"
)

func main() {
	config.LoadEnv()
	// environment := config.GetEnv("ENVIRONMENT", "development")
	// port := config.GetEnv("PORT", "8080")

	cfg := config.ConnectDataBase()

	sqlDB, _ := cfg.DB()
	defer sqlDB.Close()

	// router
	server := fiber.New()
	routes.SetupRouter(server)

	server.Listen(":" + config.AppEnv.AppPort)
}
