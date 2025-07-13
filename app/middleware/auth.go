package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2" // GANTI: Menggunakan Fiber
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware adalah middleware untuk memvalidasi token otentikasi di Fiber.
// GANTI: Tipe return diubah dari gin.HandlerFunc menjadi fiber.Handler
func AuthMiddleware() fiber.Handler {
	// GANTI: Signature fungsi diubah menjadi func(c *fiber.Ctx) error
	return func(c *fiber.Ctx) error {
		// GANTI: c.GetHeader() menjadi c.Get()
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// GANTI: c.AbortWithStatusJSON menjadi return c.Status().JSON()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header tidak ditemukan",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Format Authorization header tidak valid. Gunakan 'Bearer <token>'",
			})
		}

		tokenString := parts[1]
		jwtSecret := os.Getenv("JWT_SECRET_KEY")

		// Mem-parsing dan memvalidasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma penandatanganan adalah yang kita harapkan (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau kedaluwarsa"})
		}

		// (Opsional) Menyimpan claims dari token ke context agar bisa diakses oleh handler berikutnya
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			c.Locals("user_claims", claims)
		}
		// GANTI: c.Next() menjadi return c.Next() untuk melanjutkan ke handler berikutnya.
		return c.Next()
	}
}
