package middlewares

import (
	"encoding/base64"
	"net/http"
	"strings"

	config "go-htmlcsstoimage/configs"
	model "go-htmlcsstoimage/internal/models"

	"github.com/gofiber/fiber/v2"
)

// AuthBasic adalah middleware untuk HTTP Basic Authentication
func AuthBasic(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Unauthorized",
			"statusCode": 401,
			"message":    "Missing or invalid Authorization header",
		})
	}

	// Decode header Basic Auth
	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Unauthorized",
			"statusCode": 401,
			"message":    "Invalid Authorization header",
		})
	}

	// Memisahkan user_id dan api_key
	parts := strings.SplitN(string(decoded), ":", 2)

	var auth model.AuthBasic
	var apikey model.ApiKey
	auth.UserID, auth.APIKey = parts[0], parts[1]

	result := config.DB.Where("user_id = ? AND api_key = ?", auth.UserID, auth.APIKey).First(&apikey)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Unauthorized",
			"statusCode": 401,
			"message":    "Invalid User ID or API Key",
		})
	}

	c.Locals("UserID", auth.UserID)
	return c.Next()
}
