package gateway

import "github.com/fiap-161/tc-golunch-core-service/internal/shared/entity"

type AuthGateway interface {
	GenerateToken(userID string, userType string, additionalClaims map[string]any) (string, error)
	ValidateToken(tokenString string) (*entity.CustomClaims, error)
}
