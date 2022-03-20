package handle

import (
	"fmt"
	"gin-auth/auth/jwt"
	jwtlib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const authHeader = "Authorization"
const authTokenPrefix = "Bearer "
const ctxDataTokenKey = "token"
const ctxDataIdKey = "id"
const ctxDataUsernameKey = "username"
const ctxDataRolesKey = "roles"

func JwtAuthenticationMw(service jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader(authHeader)
		if !strings.HasPrefix(tokenHeader, authTokenPrefix) {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(tokenHeader, authTokenPrefix)
		token, err := service.VerifyToken(tokenStr)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Set(ctxDataTokenKey, token)
		fmt.Printf("%T", token.Claims)
		claims := token.Claims.(jwtlib.MapClaims)
		c.Set(ctxDataIdKey, claims[jwt.AppClaimsId])
		c.Set(ctxDataUsernameKey, claims[jwt.AppClaimsUsername])
		c.Set(ctxDataRolesKey, claims[jwt.AppClaimsRoles])
	}
}

func JwtAuthorizationHasAnyRoleMv(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		existingRoles, ok := c.Get(ctxDataTokenKey)
		if !ok {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		if !roleContainsAny(existingRoles.([]string), roles...) {
			c.Status(http.StatusForbidden)
			c.Abort()
		}
	}
}

func JwtAuthorizationHasEachRoleMv(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		existingRoles, ok := c.Get(ctxDataTokenKey)
		if !ok {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		if !roleContainsEach(existingRoles.([]string), roles...) {
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

func roleContainsEach(existingRoles []string, roles ...string) bool {
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
