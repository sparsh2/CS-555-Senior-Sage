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
	Name            string             `json:"name"`
	Msg             string             `json:"msg"`
}
