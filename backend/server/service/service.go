package service

import "server/types"

type Service struct {
}

var Svc *Service

func init() {
	Svc = &Service{}
}

func (s *Service) GetUserDetails(userId string, userToken string) (*types.GetUserDetailsResponse, error) {
	userDetails, err := StorageSvc.GetUserDetails(userId, userToken)
	return userDetails, err
}

func (s *Service) GetAccessLogs(req types.GetRequestLogsRequest) (*types.GetRequestLogsResponse, error) {
	requestLogs, err := StorageSvc.GetAccessLogs(req)
	return requestLogs, err
}
