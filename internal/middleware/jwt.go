package middleware

import (
	"beli-mang/config"
	"beli-mang/pkg/jwt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type JwtMiddleware struct {
	config config.Configs
}

// Auth protect routes
func (jwtMiddleware *JwtMiddleware) Auth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(jwtMiddleware.config.EnvConfig["JWT_SECRET"])},
		SuccessHandler: jwtMiddleware.jwtSuccess,
		ErrorHandler:   jwtMiddleware.jwtError,
	})
}

func (jwtMiddleware *JwtMiddleware) jwtSuccess(c *fiber.Ctx) error {
	claims, err := jwt.ExtractToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Locals("user_id", claims.UserID)
	return c.Next()
}

func (jwtMiddleware *JwtMiddleware) jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing or malformed JWT",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid or expired JWT",
		"data":    nil,
	})
}
