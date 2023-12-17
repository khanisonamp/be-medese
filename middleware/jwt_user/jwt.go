package jwt

import (
	"time"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func CreateTokenUser(username string) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["create_at"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token and send it as response.
	jwtToken, err := token.SignedString([]byte(viper.GetString("jwt.key")))
	if err != nil {
		return "create_token_fail", err
	}

	return jwtToken, err
}

func AuthUser() fiber.Handler {

	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(viper.GetString("jwt.key")),
		ErrorHandler: jwtErrorAuthUser,
	})
}

func jwtErrorAuthUser(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     fiber.StatusUnauthorized,
			"message":    "invalid_or_expired_jwt",
			"message_th": "token ไม่ถูกต้องหรือหมดอายุ",
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status":     fiber.StatusUnauthorized,
		"message":    "invalid_or_expired_jwt",
		"message_th": "token ไม่ถูกต้องหรือหมดอายุ",
	})
}
