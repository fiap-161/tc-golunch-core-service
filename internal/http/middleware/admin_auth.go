package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware verifica se o token pertence a um admin via Operation Service
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token with Operation Service
		if !validateAdminToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired admin token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func validateAdminToken(token string) bool {
	// Call Operation Service to validate admin token
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8083/admin/validate", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error calling operation service: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return false
	}

	// Check if the response indicates valid admin
	valid, exists := response["valid"].(bool)
	return exists && valid
}
