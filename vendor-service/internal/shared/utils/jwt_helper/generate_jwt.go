package jwt_helper

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

func GenerateVendorJwt(vendorId uuid.UUID) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	claims := jwt.MapClaims{
		"vendor_id": vendorId,
		"exp":       time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
