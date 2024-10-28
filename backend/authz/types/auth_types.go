package types

type AuthVerifyRequest struct {
	JWTToken string `json:"jwt_token"`
}

type UserDetails struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type UserLoginRequest struct {
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
}

type MongoUserDoc struct {
	UserDetails MongoUserDetails `bson:"user_details"`
	UserId      string           `bson:"_id"`
}

type MongoUserDetails struct {
	Email        string `bson:"email,omitempty"`
	PasswordHash string `bson:"password_hash,omitempty"`
}

type UserSignupRequest struct {
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
}
