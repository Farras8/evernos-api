// file: middleware/auth_middleware.go

package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Ambil header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Missing authorization header"})
	}

	// Pisahkan "Bearer" dengan tokennya
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid authorization header format"})
	}
	
	tokenString := parts[1]

	// Parse dan validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired token"})
	}

	// Ambil claims dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token claims"})
	}

	// Simpan informasi user ke context untuk digunakan di handler selanjutnya
	c.Locals("user_id", claims["user_id"])
	c.Locals("is_admin", claims["is_admin"])

	// Lanjutkan ke handler/middleware selanjutnya
	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
	// Ambil data is_admin yang sudah disimpan oleh AuthMiddleware
	isAdmin, ok := c.Locals("is_admin").(bool)

	// Jika bukan admin atau data tidak ada, kembalikan error
	if !ok || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: Admins only"})
	}

	// Lanjutkan jika user adalah admin
	return c.Next()
}