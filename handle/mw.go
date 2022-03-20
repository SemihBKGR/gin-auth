package handle

import (
	"gin-auth/auth/jwt"
	jwtlib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const authHeader = "Authorization"
const authTokenPrefix = "Bearer "
const ctxDataTokenKey = "token"

func JwtAuthenticationMw(service jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader(authHeader)
		if strings.HasPrefix(tokenHeader, authHeader) {
			token := strings.TrimPrefix(tokenHeader, authTokenPrefix)
			if token, err := service.VerifyToken(token); err == nil {
				c.Set(ctxDataTokenKey, token)
			}
		}
		c.Status(http.StatusUnauthorized)
		c.Abort()
	}
}

func JwtAuthorizationHasAnyRoleMv(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInterface, ok := c.Get(ctxDataTokenKey)
		if !ok {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		token := tokenInterface.(jwtlib.Token)
		claims := token.Claims.(jwt.AppClaims)
		if !roleContainsAny(claims.Roles, roles...) {
			c.Status(http.StatusForbidden)
			c.Abort()
		}
	}
}

func JwtAuthorizationHasAllRoleMv(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInterface, ok := c.Get(ctxDataTokenKey)
		if !ok {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		token := tokenInterface.(jwtlib.Token)
		claims := token.Claims.(jwt.AppClaims)
		if !roleContainsAll(claims.Roles, roles...) {
			c.Status(http.StatusForbidden)
			c.Abort()
		}
	}
}

func roleContainsAny(existingRoles []string, roles ...string) bool {
	for _, role := range roles {
		if roleContains(existingRoles, role) {
			return true
		}
	}
	return false
}

func roleContainsAll(existingRoles []string, roles ...string) bool {
	for _, role := range roles {
		if !roleContains(existingRoles, role) {
			return false
		}
	}
	return true
}

func roleContains(existingRoles []string, role string) bool {
	for _, existingRole := range existingRoles {
		if role == existingRole {
			return true
		}
	}
	return false
}
