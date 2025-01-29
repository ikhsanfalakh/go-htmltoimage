package routes

import (
	controller "go-htmlcsstoimage/internal/controllers"
	middleware "go-htmlcsstoimage/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(server fiber.Router) {

	server.Use(middleware.AuthBasic)
	server.Post("/v1/image/", controller.GenerateImage)
	
	server.Get("/storage/images/:filename", controller.GetImage)
}
