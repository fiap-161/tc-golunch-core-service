package gateway

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthGatewayImpl struct {
	secretKey string
}

func NewAuthGateway(secretKey string) AuthGateway {
	return &AuthGatewayImpl{
		secretKey: secretKey,
	}
}

func (a *AuthGatewayImpl) GenerateToken(userID string, userType string, additionalClaims map[string]any) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":       time.Now().Unix(),
	}

	// Add additional claims
	for key, value := range additionalClaims {
		claims[key] = value
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
