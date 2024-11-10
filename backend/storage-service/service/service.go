package service

import (
	"fmt"
	"log"
	"storage-service/storage"
	"storage-service/types"
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
		return types.ErrAccessDenied
	}
	userDetails, err := storage.StorageSvc.GetUserDoc(req.UserId)
	if err != nil {
		return fmt.Errorf("error in getting user data: %v", err)
	}
	if userDetails.ReminderDetails == nil {
		userDetails.ReminderDetails = &[]types.ReminderDetails{}
	}
	userDetails.ReminderDetails = req.Reminders
	err = storage.StorageSvc.InsertUserDoc(userDetails)
	if err != nil {
		return fmt.Errorf("error in writing user data: %v", err)
	}
	return nil
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
	err = storage.StorageSvc.InsertUserDoc(userDetails)
	if err != nil {
		return fmt.Errorf("error in writing user data: %v", err)
	}
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
	err = storage.StorageSvc.InsertUserDoc(userDetails)
	if err != nil {
		return fmt.Errorf("error in writing user data: %v", err)
	}
	return nil
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
	resp.Msg = "data fetched successfully"
	return resp, nil
}
