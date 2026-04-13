package middleware

import (
	"strings"

	"github.com/falaqmsi/go-example/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/MicahParks/keyfunc/v3"
)

const (
	CtxUserIDKey = "user_id"
	CtxEmailKey  = "email"
	CtxRolesKey  = "roles"
)

// Auth is a middleware that verifies JWT tokens using Keycloak's JWKS.
func Auth(jwks keyfunc.Keyfunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is missing")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse the JWT using the JWKS keyfunc payload
		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "Failed to parse token claims")
			c.Abort()
			return
		}

		// Extract standard JWT claims and Keycloak specific claims
		sub, _ := claims["sub"].(string) // user_id
		email, _ := claims["email"].(string)

		var roles []string
		
		// Typically keycloak places roles in realm_access.roles
		if realmAccessInter, ok := claims["realm_access"]; ok {
			if realmAccess, ok := realmAccessInter.(map[string]interface{}); ok {
				if rolesInter, ok := realmAccess["roles"].([]interface{}); ok {
					for _, r := range rolesInter {
						if roleStr, ok := r.(string); ok {
							roles = append(roles, roleStr)
						}
					}
				}
			}
		}

		// Save into context for downstream handlers
		c.Set(CtxUserIDKey, sub)
		c.Set(CtxEmailKey, email)
		c.Set(CtxRolesKey, roles)

		c.Next()
	}
}
