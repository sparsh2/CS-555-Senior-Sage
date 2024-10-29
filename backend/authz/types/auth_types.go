package types

type AuthVerifyRequest struct {
	JWTToken string `json:"jwt_token"`
}

type UserDetails struct {
	UserId    string `json:"user_id"`
	UserEmail string `json:"username"`
}

type UserLoginRequest struct {
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
}

type MongoUserDoc struct {
	UserDetails MongoUserDetails `bson:"user_details,omitempty"`
	UserId      string           `bson:"_id,omitempty"`
}

type MongoUserDetails struct {
	Email        string `bson:"email,omitempty"`
	PasswordHash string `bson:"password_hash,omitempty"`
}

type UserSignupRequest struct {
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
}
