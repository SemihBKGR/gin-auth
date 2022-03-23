package jwt

import (
	"fmt"
	"gin-auth/persist"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtService interface {
	GenerateToken(user *persist.User) string
	VerifyToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	Secret []byte
	Issuer string
}

const AppClaimsUsername = "Username"
const AppClaimsRoles = "Roles"

type AppClaims struct {
	*jwt.StandardClaims
	ID       uint
	Username string
	Roles    []string
}

func (s *jwtService) GenerateToken(user *persist.User) string {
	claims := &AppClaims{
		StandardClaims: &jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    s.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
		ID:       user.ID,
		Username: user.Username,
		Roles:    rolesToString(user.Roles),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.Secret)
	if err != nil {
		panic(err)
	}
	return tokenStr
}

func (s *jwtService) VerifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("invalid token, alg: %s", token.Header["alg"])
		}
		return s.Secret, nil
	})
}

func NewJwtService(secret, issuer string) JwtService {
	return &jwtService{
		Secret: []byte(secret),
		Issuer: issuer,
	}
}

func rolesToString(roles []persist.Role) []string {
	rolesStr := make([]string, len(roles))
	for i, role := range roles {
		rolesStr[i] = role.Name
	}
	return rolesStr
}
