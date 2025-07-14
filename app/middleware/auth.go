package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware adalah middleware untuk memvalidasi token otentikasi di Fiber.
// Middleware ini akan mencari token dari:
// 1. Header `Authorization: Bearer <token>`
// 2. Query parameter `?token=<token>`
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		// Coba dapatkan token dari header 'Authorization'
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Jika tidak ada di header, coba dapatkan dari query parameter 'token'
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		// Jika token masih tidak ditemukan, kembalikan error
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token otentikasi tidak ditemukan di header atau query parameter",
			})
		}

		jwtSecret := os.Getenv("JWT_SECRET_KEY")

		// Mem-parsing dan memvalidasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma penandatanganan adalah yang kita harapkan (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode penandatanganan tidak terduga: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau kedaluwarsa"})
		}

		// (Opsional) Menyimpan claims dari token ke context agar bisa diakses oleh handler berikutnya
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_claims", claims)
		}
		return c.Next()
	}
}
