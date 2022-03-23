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
const ctxDataClaimsKey = "claims"
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
		claims := token.Claims.(jwtlib.MapClaims)
		c.Set(ctxDataClaimsKey, claims)
		c.Set(ctxDataUsernameKey, claims[jwt.AppClaimsUsername])
		c.Set(ctxDataRolesKey, claims[jwt.AppClaimsRoles])
	}
}

func JwtAuthorizationHasAnyRoleMv(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		existingRoles, ok := ExtractRolesContextData(c)
		if !ok {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		if !roleContainsAny(existingRoles, roles...) {
			c.Status(http.StatusForbidden)
			c.Abort()
		}
	}
}

func JwtAuthorizationHasEachRoleMv(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		existingRoles, ok := ExtractRolesContextData(c)
		if !ok {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		if !roleContainsEach(existingRoles, roles...) {
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

func ExtractUsernameContextData(c *gin.Context) (string, bool) {
	usernameData, ok := c.Get(ctxDataUsernameKey)
	if !ok {
		return "", false
	}
	username, ok := usernameData.(string)
	if !ok {
		return "", false
	}
	return username, true
}

func ExtractRolesContextData(c *gin.Context) ([]string, bool) {
	rolesData, ok := c.Get(ctxDataRolesKey)
	if !ok {
		return nil, false
	}
	rolesInterface, ok := rolesData.([]interface{})
	if !ok {
		return nil, false
	}
	roles := make([]string, len(rolesInterface))
	for i, role := range rolesInterface {
		roles[i] = role.(string)
	}
	return roles, true
}
