package types

import (
	"fmt"
)

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

type RequestAccessResponse struct {
	AccessRequest bool   `json:"access_request"`
	Message       string `json:"message,omitempty"`
}

type ResourceType string

const (
	RESOURCE_USER_DETAILS      ResourceType = "USER_DETAILS"
	RESOURCE_USER_PREFERENCES  ResourceType = "USER_PREFERENCES"
	RESOURCE_USER_REMINDERS    ResourceType = "USER_REMINDERS"
	RESOURCE_USER_CHAT_HISTORY ResourceType = "USER_CHAT_HISTORY"
	RESOURCE_RPM_READINGS      ResourceType = "RPM_READINGS"
	RESOURCE_UNKNOWN           ResourceType = "UNKNOWN"
)

var ResourceList = []ResourceType{
	RESOURCE_USER_CHAT_HISTORY,
	RESOURCE_USER_PREFERENCES,
	RESOURCE_USER_REMINDERS,
	RESOURCE_RPM_READINGS,
}

const (
	MONGO_USER_DETAILS_FIELD     = "user_details"
	MONGO_USER_PREFERENCES_FIELD = "user_preferences"
	MONGO_USER_REMINDERS_FIELD   = "user_reminders"
)

var ResourceToMongoField map[ResourceType]string = map[ResourceType]string{
	RESOURCE_USER_DETAILS:     MONGO_USER_DETAILS_FIELD,
	RESOURCE_USER_PREFERENCES: MONGO_USER_PREFERENCES_FIELD,
	RESOURCE_USER_REMINDERS:   MONGO_USER_REMINDERS_FIELD,
}

var ResourceIdToString map[ResourceType]string = map[ResourceType]string{
	RESOURCE_USER_DETAILS:     "UserDetailsResource",
	RESOURCE_USER_PREFERENCES: "UserPreferenceResource",
	RESOURCE_USER_REMINDERS:   "UserRemindersResource",
}

// type MongoUserDoc struct {
// 	UserDetails MongoUserDetails  `bson:"user_details,omitempty"`
// 	UserId      string            `bson:"_id,omitempty"`
// 	Data        MongoResourceData `bson:"data,omitempty"`
// }

type MongoResourceData struct {
	UserDetails     string `bson:"user_details,omitempty"`
	UserPreferences string `bson:"user_preferences,omitempty"`
	UserReminders   string `bson:"user_reminders,omitempty"`
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

type WriteDataResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

type ResourceData struct {
	UserDetails     string `json:"user_details,omitempty"`
	UserPreferences string `json:"user_preferences,omitempty"`
	UserReminders   string `json:"user_reminders,omitempty"`
}

type WriteDataRequest struct {
	UserId         string                  `json:"user_id"`
	RequesterToken string                  `json:"requester_id"`
	Data           map[ResourceType]string `json:"data"`
}

var ErrInvalidToken = fmt.Errorf("invalid token")
var ErrAccessDenied error = fmt.Errorf("access denied")
