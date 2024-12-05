package service

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"log"
	"storage-service/storage"
	"storage-service/types"
	"strings"
)

type Service struct {
}

var Svc *Service

func init() {
	Svc = &Service{}
}

func (as *Service) WriteReminders(req *types.WriteRemindersRequest) error {
	granted, _, err := AuthzClientSvc.VerifyAccessRequest(
		req.RequesterToken,
		req.UserId,
		[]types.ResourceType{types.RESOURCE_USER_REMINDERS},
	)
	if err != nil {
		return fmt.Errorf("error in verifying access request: %v", err)
	}
	if !granted {
		logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_REMINDERS}, granted)
		return types.ErrAccessDenied
	}
	userDetails, err := storage.StorageSvc.GetUserDoc(req.UserId)
	if err != nil {
		return fmt.Errorf("error in getting user data: %v", err)
	}
	if userDetails.ReminderDetails == nil {
		userDetails.ReminderDetails = &[]types.ReminderDetails{}
	}
	userDetails.ReminderDetails = &[]types.ReminderDetails{}
	*userDetails.ReminderDetails = append(*userDetails.ReminderDetails, *req.Reminders...)
	err = storage.StorageSvc.UpdateUserDoc(userDetails)
	if err != nil {
		return fmt.Errorf("error in writing user data: %v", err)
	}
	logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_REMINDERS}, granted)
	return nil
}

func logAccess(requesterToken, userId, operation string, resources []types.ResourceType, granted bool) {
	email, err := getEmailFromJWT(requesterToken)
	if err != nil {
		log.Printf("could not log access request: error in getting email from jwt: %v\n", err)
		return
	}
	for _, res := range resources {
		err = storage.StorageSvc.LogAccess(email, userId, operation, res, granted)
		if err != nil {
			log.Printf("error in logging access: %v\n", err)
			return
		}
		log.Printf("access logged successfully\n")
	}
}

func getEmailFromJWT(jwtToken string) (string, error) {
	parts := strings.Split(jwtToken, ".")

	// Decode the payload part (index 1)
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("Error decoding payload:", err)
		return "", err
	}

	// Parse the JSON payload
	var claims map[string]interface{}
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		fmt.Println("Error parsing claims:", err)
		return "", err
	}

	email, ok := claims["user_id"].(string)
	if !ok {
			fmt.Println("Email claim not found")
			return "", fmt.Errorf("Email claim not found")
	}

	return email, nil
}

func (as *Service) WriteResponses(req *types.WriteResponsesRequest) error {
	granted, _, err := AuthzClientSvc.VerifyAccessRequest(
		req.RequesterToken,
		req.UserId,
		[]types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY},
	)
	if err != nil {
		return fmt.Errorf("error in verifying access request: %v", err)
	}
	if !granted {
		logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY}, granted)
		return types.ErrAccessDenied
	}
	userDetails, err := storage.StorageSvc.GetUserDoc(req.UserId)
	if err != nil {
		return fmt.Errorf("error in getting user data: %v", err)
	}
	if userDetails.ChatHistory == nil {
		userDetails.ChatHistory = &[]types.ChatSession{}
	}
	if len(*req.Responses) < len(*userDetails.ChatHistory) {
		log.Printf("warning: new responses are less than the existing responses for user %s\n", req.UserId)
	}
	userDetails.QuestionResponses = req.Responses
	err = storage.StorageSvc.UpdateUserDoc(userDetails)
	if err != nil {
		logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY}, granted)
		return fmt.Errorf("error in writing user data: %v", err)
	}
	return nil
}

func (as *Service) WriteChatHistory(req *types.WriteChatHistoryRequest) error {
	granted, _, err := AuthzClientSvc.VerifyAccessRequest(
		req.RequesterToken,
		req.UserId,
		[]types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY},
	)
	if err != nil {
		return fmt.Errorf("error in verifying access request: %v", err)
	}
	if !granted {
		logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY}, granted)
		return types.ErrAccessDenied
	}
	userDetails, err := storage.StorageSvc.GetUserDoc(req.UserId)
	if err != nil {
		return fmt.Errorf("error in getting user data: %v", err)
	}
	if userDetails.ChatHistory == nil {
		userDetails.ChatHistory = &[]types.ChatSession{}
	}
	*userDetails.ChatHistory = append(*userDetails.ChatHistory, req.ChatHistory)
	err = storage.StorageSvc.UpdateUserDoc(userDetails)
	if err != nil {
		return fmt.Errorf("error in writing user data: %v", err)
	}
	logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY}, granted)
	return nil
}

func (as *Service) WriteQuestionCounter(req *types.WriteQuestionCounterRequest) error {
	granted, _, err := AuthzClientSvc.VerifyAccessRequest(
		req.RequesterToken,
		req.UserId,
		[]types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY},
	)
	if err != nil {
		return fmt.Errorf("error in verifying access request: %v", err)
	}
	if !granted {
		logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY}, granted)
		return types.ErrAccessDenied
	}
	userDetails, err := storage.StorageSvc.GetUserDoc(req.UserId)
	if err != nil {
		return fmt.Errorf("error in getting user data: %v", err)
	}
	userDetails.QuestionCounts = req.QuestionCounts
	err = storage.StorageSvc.UpdateUserDoc(userDetails)
	if err != nil {
		return fmt.Errorf("error in writing user data: %v", err)
	}
	logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_CHAT_HISTORY}, granted)
	return nil
}

func (as *Service) WritePreferences(req *types.WritePreferencesRequest) error {
	granted, _, err := AuthzClientSvc.VerifyAccessRequest(
		req.RequesterToken,
		req.UserId,
		[]types.ResourceType{types.RESOURCE_USER_PREFERENCES},
	)
	if err != nil {
		return fmt.Errorf("error in verifying access request: %v", err)
	}
	if !granted {
		logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_PREFERENCES}, granted)
		return types.ErrAccessDenied
	}
	userDetails, err := storage.StorageSvc.GetUserDoc(req.UserId)
	if err != nil {
		return fmt.Errorf("error in getting user data: %v", err)
	}
	if userDetails.Preferences == nil {
		userDetails.Preferences = &[]string{}
	}

	if len(*req.Preferences) < len(*userDetails.Preferences) {
		log.Printf("warning: new preferences are less than the existing preferences for user %s\n", req.UserId)
	}
	userDetails.Preferences = req.Preferences
	err = storage.StorageSvc.UpdateUserDoc(userDetails)
	if err != nil {
		return fmt.Errorf("error in writing user data: %v", err)
	}
	logAccess(req.RequesterToken, req.UserId, types.OPERATION_WRITE, []types.ResourceType{types.RESOURCE_USER_PREFERENCES}, granted)
	return nil
}

func (as *Service) GetAccessLogs(userEmail string) ([]string, error) {
	accessLogs, err := storage.StorageSvc.GetAccessLogs(userEmail)
	if err != nil {
		return []string{}, fmt.Errorf("error in getting access logs: %v", err)
	}
	logsString := []string{}
	for _, log := range accessLogs {
		deniedOrGranted := "denied"
		if log.Granted {
			deniedOrGranted = "granted"
		}
		logsString = append(logsString, fmt.Sprintf("User %s requested %s access on %s.%s at %s; request was %s", log.RequesterEmail, log.Operation, log.UserEmail, log.Resource, log.TimeStamp, deniedOrGranted))
	}
	return logsString, nil
}

func (as *Service) GetData(getDataReq *types.GetDataRequest) (*types.GetDataResponse, error) {
	granted, msg, err := AuthzClientSvc.VerifyAccessRequest(
		getDataReq.RequesterToken,
		getDataReq.UserId,
		[]types.ResourceType{
			types.RESOURCE_USER_PREFERENCES,
			types.RESOURCE_USER_REMINDERS,
			types.RESOURCE_USER_CHAT_HISTORY,
			types.RESOURCE_RPM_READINGS,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error in verifying access request: %v", err)
	}
	if !granted {
		logAccess(getDataReq.RequesterToken, getDataReq.UserId, types.OPERATION_READ, []types.ResourceType{
			types.RESOURCE_USER_PREFERENCES,
			types.RESOURCE_USER_REMINDERS,
			types.RESOURCE_USER_CHAT_HISTORY,
			types.RESOURCE_RPM_READINGS,
		}, granted)
		return &types.GetDataResponse{Msg: msg}, types.ErrAccessDenied
	}
	resp := &types.GetDataResponse{}
	userDetails, err := storage.StorageSvc.GetUserDoc(getDataReq.UserId)
	if err != nil {
		return nil, fmt.Errorf("error in getting user data: %v", err)
	}
	resp.ChatHistory = userDetails.ChatHistory
	resp.Preferences = userDetails.Preferences
	resp.ReminderDetails = userDetails.ReminderDetails
	resp.RPMReadings = userDetails.RPMReadings
	resp.VoiceSelection = userDetails.VoiceSelection
	resp.QuestionResponses = userDetails.QuestionResponses
	resp.QuestionCounts = userDetails.QuestionCounts
	resp.Name = userDetails.Name
	resp.Msg = "data fetched successfully"
	logAccess(getDataReq.RequesterToken, getDataReq.UserId, types.OPERATION_READ, []types.ResourceType{
		types.RESOURCE_USER_PREFERENCES,
		types.RESOURCE_USER_REMINDERS,
		types.RESOURCE_USER_CHAT_HISTORY,
		types.RESOURCE_RPM_READINGS,
	}, granted)
	return resp, nil
}
