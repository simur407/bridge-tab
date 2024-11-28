package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JwtGuard() fiber.Handler {
	return jwtware
}
