package service

import (
	"authz/config"
	"authz/types"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationService struct {
}

var AuthService *AuthenticationService

func init() {
	AuthService = &AuthenticationService{}
}

func (as *AuthenticationService) GenerateToken(userDetails *types.UserDetails) (string, error) {
	signingKey := []byte(config.Configs.AuthSecretKey)
	claims := &jwt.RegisteredClaims{
		Issuer:    "sage-server",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(9999 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	return ss, err
}

func (as *AuthenticationService) VerifyToken(tokenString string) (bool, error) {
	token, err  := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {return []byte(config.Configs.AuthSecretKey), nil})
	if err != nil {
		return false, fmt.Errorf("unable to verify the token: %v", err)
	}
	if !token.Valid {
		return false, nil
	}
	return true, nil
}
