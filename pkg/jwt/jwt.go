package jwt

import (
	"beli-mang/config"
	"beli-mang/pkg/str"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	secretKey = config.GetConfig("JWT_SECRET")
	expSecret = str.StringToInt(config.GetConfig("JWT_EXP_SECRET"))
)

type JWTCustomClaim struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userId *uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expSecret) * time.Hour).Unix()

	claims := JWTCustomClaim{
		UserID: *userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime, // 8 hour
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

// ExtractToken extracts token from echo context
func ExtractToken(c *fiber.Ctx) (*JWTCustomClaim, error) {
	tokenFromHeader := c.Get(fiber.HeaderAuthorization)
	if tokenFromHeader == "" {
		return nil, fiber.ErrUnauthorized
	}

	// Remove "Bearer " prefix
	tokenParts := strings.Split(tokenFromHeader, " ")
	if len(tokenParts) != 2 {
		return nil, fiber.ErrUnauthorized
	}
	sanitizedToken := tokenParts[1]

	token, err := jwt.Parse(sanitizedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	if !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	customClaims := &JWTCustomClaim{
		UserID: userID,
	}

	return customClaims, nil
}
