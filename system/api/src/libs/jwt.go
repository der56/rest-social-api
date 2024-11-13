package libs

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID string `json:"sub"`
	jwt.StandardClaims
}

func GenerateJWT(userID interface{}) (string, error) {
	secret := os.Getenv("SECRET_JWT")
	if secret == "" {
		return "", errors.New("La variable de entorno SECRET_JWT no está configurada")
	}

	var userIDString string

	switch v := userID.(type) {
	case string:
		userIDString = v
	case uuid.UUID:
		userIDString = v.String()
	default:
		return "", errors.New("Tipo de userID no admitido")
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	if expirationTime.IsZero() {
		return "", errors.New("Error obteniendo la hora de expiración")
	}

	claims := CustomClaims{
		UserID: userIDString,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
