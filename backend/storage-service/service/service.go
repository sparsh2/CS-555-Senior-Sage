package service

import (
	"fmt"
	"storage-service/storage"
	"storage-service/types"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
}

var Svc *Service

func init() {
	Svc = &Service{}
}

type CustomRegisteredClaims struct {
	jwt.RegisteredClaims
	UserId    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}

func (as *Service) WriteData(writeDataReq *types.WriteDataRequest) (*types.WriteDataResponse, error) {
	resources := []types.ResourceType{}
	for resource := range writeDataReq.Data {
		switch resource {
		case types.RESOURCE_USER_DETAILS:
			resources = append(resources, types.RESOURCE_USER_DETAILS)
		case types.RESOURCE_USER_PREFERENCES:
			resources = append(resources, types.RESOURCE_USER_PREFERENCES)
		case types.RESOURCE_USER_REMINDERS:
			resources = append(resources, types.RESOURCE_USER_REMINDERS)
		}
	}
	granted, msg, err := AuthzClientSvc.VerifyAccessRequest(writeDataReq.RequesterToken, writeDataReq.UserId, resources)
	if err != nil {
		return nil, fmt.Errorf("error in verifying access request: %v", err)
	}
	if !granted {
		return &types.WriteDataResponse{Success: false, Msg: msg}, types.ErrAccessDenied
	}
	err = storage.StorageSvc.WriteUserData(writeDataReq.UserId, writeDataReq.Data, resources)
	if err != nil {
		return nil, fmt.Errorf("error in writing user data: %v", err)
	}
	return &types.WriteDataResponse{}, nil
}

func (as *Service) GetData(getDataReq *types.GetDataRequest) (*types.GetDataResponse, error) {
	granted, msg, err := AuthzClientSvc.VerifyAccessRequest(getDataReq.RequesterToken, getDataReq.UserId, getDataReq.Resources)
	if err != nil {
		return nil, fmt.Errorf("error in verifying access request: %v", err)
	}
	if !granted {
		return &types.GetDataResponse{Msg: msg}, types.ErrAccessDenied
	}
	resp := &types.GetDataResponse{
		Data: make(map[types.ResourceType]string),
	}
	err = storage.StorageSvc.GetUserData(getDataReq.UserId, getDataReq.Resources, resp)
	if err != nil {
		return nil, fmt.Errorf("error in getting user data: %v", err)
	}
	return resp, nil
}
