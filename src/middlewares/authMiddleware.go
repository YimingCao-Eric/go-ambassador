package middlewares

import (
	"golangProject/util"

	"github.com/gofiber/fiber/v3"
)

// IsAuthenticated middleware checks if the request has a valid JWT token
// This protects routes that require user authentication
// Usage: app.Use(middlewares.IsAuthenticated) or app.Get("/protected", middleware, handler)
func IsAuthenticated(c fiber.Ctx) error {
	// Extract JWT token from the "jwt" cookie
	cookie := c.Cookies("jwt")

	// Validate the token using the utility function
	if _, err := util.ParseJWT(cookie); err != nil {
		c.Status(fiber.StatusUnauthorized) // Set HTTP status to 401 Unauthorized
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	// If token is valid, proceed to the next handler in the chain
	return c.Next()
}
