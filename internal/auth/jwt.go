package auth

import (
	"time"
	"os" 

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateToken membuat token JWT baru untuk user
func GenerateToken(userID string) (string, error) {
	// Set claims (payload) untuk token
	claims := jwt.MapClaims{
		"sub": userID, // Subject (ID Pengguna)
		"iat": time.Now().Unix(), // Issued At
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expiration time (24 jam)
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}