package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuthenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing Authorization header"})
		}

		tokenString := authHeader[len("Bearer "):]
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token format"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		if authVal, exists := claims["authenticated"]; !exists || authVal != true {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authentication required"})
		}

		return c.Next()
	}
}
