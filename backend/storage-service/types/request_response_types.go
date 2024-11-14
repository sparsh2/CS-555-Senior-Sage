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
	ReminderDetails   *[]ReminderDetails      `json:"reminder_details"`
	RPMReadings       *[]RPMReading           `json:"rpm_readings"`
	Preferences       *[]string               `json:"preferences"`
	ChatHistory       *[]ChatSession          `json:"chat_history"`
	QuestionResponses *[]QuestionResponse     `json:"question_responses"`
	VoiceSelection    string                  `json:"voice_selection"`
	Name              string                  `json:"name"`
	Msg               string                  `json:"msg"`
	QuestionCounts    map[int]QuestionCounter `json:"question_counts"`
}

type QuestionCounter struct {
	Counter   bool   `json:"counter"`
	Frequency int    `json:"frequency"`
	AskedDate string `json:"asked_date"`
	CurrDate  string `json:"curr_date"`
	Diff      int    `json:"diff"`
}

type WriteQuestionCounterRequest struct {
	RequesterToken string                  `json:"requester_token"`
	UserId         string                  `json:"user_id"`
	QuestionCounts map[int]QuestionCounter `json:"question_counts"`
}

type WriteChatHistoryRequest struct {
	RequesterToken string       `json:"requester_token"`
	UserId         string       `json:"user_id"`
	ChatHistory    *[]ChatSession `json:"chat_history"`
}
