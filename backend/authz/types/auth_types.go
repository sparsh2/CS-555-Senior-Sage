package types

type AuthVerifyRequest struct {
	JWTToken string `json:"jwt_token"`
}

type UserDetails struct {
	UserId string `json:"user_id"`
	Username string `json:"username"`
}