package types

type WritePreferencesRequest struct {
	RequesterToken string    `json:"requester_token"`
	UserId         string    `json:"user_id"`
	Preferences    *[]string `json:"preferences"`
}

type WriteResponsesRequest struct {
	RequesterToken string              `json:"requester_token"`
	UserId         string              `json:"user_id"`
	Responses      *[]QuestionResponse `json:"responses"`
}

type WriteRemindersRequest struct {
	RequesterToken string             `json:"requester_token"`
	UserId         string             `json:"user_id"`
	Reminders      *[]ReminderDetails `json:"reminders"`
}

type GetDataRequest struct {
	UserId         string `json:"user_id"`
	RequesterToken string `json:"requester_token"`
}

type GetDataResponse struct {
	ReminderDetails *[]ReminderDetails `json:"reminder_details"`
	RPMReadings     *[]RPMReading      `json:"rpm_readings"`
	Preferences     *[]string          `json:"preferences"`
	ChatHistory     *[]ChatSession     `json:"chat_history"`
	VoiceSelection  string             `json:"voice_selection"`
	Msg             string             `json:"msg"`
}

// Request types
type AuthVerifyRequest struct {
	JWTToken string `json:"jwt_token"`
}

type UserLoginRequest struct {
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
}

type UserLoginResponse struct {
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type UserSignupRequest struct {
	VoiceSelection string `json:"voice_selection"`
	Name           string `json:"name"`
	UserEmail      string `json:"email"`
	UserPassword   string `json:"password"`
}

type RequestAccessRequest struct {
	RequesterToken string         `json:"requester_id"`
	UserId      string         `json:"user_id"`
	Resources   []ResourceType `json:"resources"`
}

type RequestAccessResponse struct {
	AccessRequest bool   `json:"access_request"`
	Message       string `json:"message,omitempty"`
}
