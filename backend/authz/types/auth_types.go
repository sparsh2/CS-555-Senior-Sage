package types

import "fmt"

// Request types
type AuthVerifyRequest struct {
	JWTToken string `json:"jwt_token"`
}

type UserLoginRequest struct {
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
}

type UserSignupRequest struct {
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
}

type RequestAccessRequest struct {
	RequesterToken string         `json:"requester_token"`
	UserId         string         `json:"user_id"`
	Resources      []ResourceType `json:"resources"`
}

type ResourceType int

const (
	UserDetailsResource ResourceType = iota
	UserPreferenceResource
	UserRemindersResource
)

var ResourceIdToString map[ResourceType]string = map[ResourceType]string{
	UserDetailsResource:    "UserDetailsResource",
	UserPreferenceResource: "UserPreferenceResource",
	UserRemindersResource:  "UserRemindersResource",
}

type UserDetails struct {
	UserId    string `json:"user_id"`
	UserEmail string `json:"username"`
}

type MongoUserDoc struct {
	UserDetails MongoUserDetails `bson:"user_details,omitempty"`
	UserId      string           `bson:"_id,omitempty"`
}

type MongoUserDetails struct {
	Email        string `bson:"email,omitempty"`
	PasswordHash string `bson:"password_hash,omitempty"`
}

type MongoAclsDoc struct {
	AclId  string                    `bson:"_id,omitempty"`
	UserId string                    `bson:"uid,omitempty"`
	Acls   map[ResourceType][]string `bson:"acls,omitempty"`
}

var ErrInvalidToken = fmt.Errorf("invalid token")
