package service

import (
	"authz/config"
	"authz/storage"
	"authz/types"
	"crypto/sha256"
	"encoding/hex"
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

func (as *AuthenticationService) Login(loginReq *types.UserLoginRequest) (string, error) {
	hash, err := storage.StorageSvc.GetUserHash(loginReq.UserEmail)
	if err != nil {
		return "", fmt.Errorf("error getting password hash: %v", err)
	}
	h := sha256.New()
	_, err  = h.Write([]byte(loginReq.UserPassword))
	if err != nil {
		return "", fmt.Errorf("error getting hash: %v", err)
	}
	if hash != hex.EncodeToString(h.Sum(nil)) {
		return "", fmt.Errorf("incorrect username or password")
	}
	userId, err := storage.StorageSvc.GetUserId(loginReq.UserEmail)
	if err != nil {
		return "", fmt.Errorf("error getting user id: %v", err)
	}

	return as.GenerateToken(&types.UserDetails{
		UserId: userId,
		Username: loginReq.UserEmail,
	})
}

func (as *AuthenticationService) VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { return []byte(config.Configs.AuthSecretKey), nil })
	if err != nil {
		return false, fmt.Errorf("unable to verify the token: %v", err)
	}
	if !token.Valid {
		return false, nil
	}
	return true, nil
}
