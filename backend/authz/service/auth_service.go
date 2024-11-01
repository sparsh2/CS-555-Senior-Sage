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
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthenticationService struct {
}

var AuthService *AuthenticationService

func init() {
	AuthService = &AuthenticationService{}
}

type CustomRegisteredClaims struct {
	jwt.RegisteredClaims
	UserId    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}

func (as *AuthenticationService) GenerateToken(userDetails *types.UserDetails) (string, error) {
	signingKey := []byte(config.Configs.AuthSecretKey)
	claims := &CustomRegisteredClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "sage-server",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(9999 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserId:    userDetails.UserId,
		UserEmail: userDetails.UserEmail,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	return ss, err
}

func (as *AuthenticationService) Signup(signupReq *types.UserSignupRequest) (string, error) {
	_, err := storage.StorageSvc.GetUserId(signupReq.UserEmail)
	if err == nil {
		return "", fmt.Errorf("user already exists")
	} else if err != mongo.ErrNoDocuments {
		return "", fmt.Errorf("error getting user id: %v", err)
	}
	h := sha256.New()
	_, err = h.Write([]byte(signupReq.UserPassword))
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	passwordHash := hex.EncodeToString(h.Sum(nil))
	userDoc := &types.MongoUserDoc{
		UserDetails: types.MongoUserDetails{
			Email:        signupReq.UserEmail,
			PasswordHash: passwordHash,
		},
		Data: map[string]string{
			types.ResourceToMongoField[types.RESOURCE_USER_DETAILS]:     "",
			types.ResourceToMongoField[types.RESOURCE_USER_PREFERENCES]: "",
			types.ResourceToMongoField[types.RESOURCE_USER_REMINDERS]:   "",
		},
	}
	err = storage.StorageSvc.InsertUserDoc(userDoc)
	if err != nil {
		return "", fmt.Errorf("error inserting user: %v", err)
	}
	uid, err := storage.StorageSvc.GetUserId(signupReq.UserEmail)
	if err != nil {
		return "", fmt.Errorf("error getting user id: %v", err)
	}
	llmuid, err := storage.StorageSvc.GetUserId(config.Configs.LLMUserEmail)
	if err != nil {
		return "", fmt.Errorf("error getting llm user id: %v", err)
	}
	aclDoc := &types.MongoAclsDoc{
		UserId: uid,
		Acls: map[types.ResourceType][]string{
			types.RESOURCE_USER_DETAILS:    {llmuid, uid},
			types.RESOURCE_USER_PREFERENCES: {llmuid},
			types.RESOURCE_USER_REMINDERS:  {llmuid, uid},
		},
	}
	err = storage.StorageSvc.InsertAclDoc(aclDoc)
	if err != nil {
		return "", fmt.Errorf("error populating ACLs: %v", err)
	}
	token, err := as.GenerateToken(&types.UserDetails{
		UserId:    uid,
		UserEmail: signupReq.UserEmail,
	})
	if err != nil {
		return "", fmt.Errorf("error generating jwt token: %v", err)
	}
	return token, nil

}

func (as *AuthenticationService) Login(loginReq *types.UserLoginRequest) (string, error) {
	hash, err := storage.StorageSvc.GetUserHash(loginReq.UserEmail)
	if err != nil {
		return "", fmt.Errorf("error getting password hash: %v", err)
	}
	h := sha256.New()
	_, err = h.Write([]byte(loginReq.UserPassword))
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
		UserId:    userId,
		UserEmail: loginReq.UserEmail,
	})
}

func (as *AuthenticationService) VerifyToken(tokenString string) (*CustomRegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomRegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte(config.Configs.AuthSecretKey), nil })
	if err != nil {
		return nil, fmt.Errorf("unable to verify the token: %v", err)
	}
	if !token.Valid {
		return nil, types.ErrInvalidToken
	}
	return token.Claims.(*CustomRegisteredClaims), nil
}

func (as *AuthenticationService) RequestAccess(requestAccessReq *types.RequestAccessRequest) (bool, error) {
	claims, err := as.VerifyToken(requestAccessReq.RequesterToken)
	if err != nil {
		return false, err
	}
	requesterId := claims.UserId
	docId := requestAccessReq.UserId
	acl, err := storage.StorageSvc.GetAclDoc(docId)
	if err != nil {
		return false, fmt.Errorf("error getting acls: %v", err)
	}
	granted := true
	for _, res := range requestAccessReq.Resources {
		ok := false
		_, ok1 := acl.Acls[res]
		if !ok1 {
			return false, fmt.Errorf("unknown resource type %v", res)
		}
		for _, uid := range acl.Acls[res] {
			if uid == requesterId {
				ok = true
				break
			}
		}
		if !ok {
			granted = false
			break
		}
	}

	if !granted {
		return false, nil
	}
	return true, nil
}
