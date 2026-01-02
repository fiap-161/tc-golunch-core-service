package middleware

import (
	"net/http"
	"strings"

	"github.com/fiap-161/tc-golunch-core-service/internal/shared/httpclient"
	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware validates admin JWT tokens via serverless Lambda
func AdminAuthMiddleware() gin.HandlerFunc {
	authClient := httpclient.NewServerlessAuthClient()

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

		// Validate admin token with serverless Lambda
		authResp, err := authClient.ValidateAdminToken(token)
		if err != nil || !authResp.Valid {
			errorMsg := "Invalid or expired admin token"
			if err != nil {
				errorMsg = err.Error()
			} else if authResp.Error != "" {
				errorMsg = authResp.Error
			}
			
			c.JSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
			c.Abort()
			return
		}

		// Set admin information in context
		c.Set("admin_id", authResp.UserID)
		c.Set("admin_role", authResp.Role)
		c.Set("admin_claims", authResp.Claims)
		
		c.Next()
	}
}
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
