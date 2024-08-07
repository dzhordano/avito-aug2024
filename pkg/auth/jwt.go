package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokensManager interface {
	GenerateJWT(userType string) (string, error)
	Parse(token string) (string, error)
}

type JWTManager struct {
	secretKey string
	tokenTTL  time.Duration
}

func NewJWTManager(secretKey string, tokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
	}
}

// GenerateJWT generates and returns JWT token from user type and sets expiration date.
func (m *JWTManager) GenerateJWT(userType string) (string, error) {
	claims := jwt.MapClaims{
		"userType": userType,
		"exp":      time.Now().Add(m.tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secretKey))
}

// Parse Returns user type from token if there is one.
func (m *JWTManager) Parse(inpToken string) (string, error) {
	token, err := jwt.Parse(inpToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["userType"].(string), nil
}
