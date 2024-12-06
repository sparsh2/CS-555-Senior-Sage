package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/config"
	"server/types"
	"strconv"
	"strings"
)

type StorageServiceClient struct {
}

var StorageSvc *StorageServiceClient

func init() {
	StorageSvc = &StorageServiceClient{}
}

func (s *StorageServiceClient) GetAccessLogs(req types.GetRequestLogsRequest) (*types.GetRequestLogsResponse, error) {
	host := config.Configs.StorageSvcConfig.Host
	port := strconv.Itoa(config.Configs.StorageSvcConfig.Port)
	// connect to storage service
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	hreq, err := http.NewRequest("GET", "http://"+host+":"+port+"/access-logs", strings.NewReader(string(body)))
	if err != nil {
		return	nil, err
	}
	res, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("error connecting to storage service: %v", err)
	}
	defer res.Body.Close()
	requestLogs := &types.GetRequestLogsResponse{}
	err = json.NewDecoder(res.Body).Decode(requestLogs)
	if err != nil {
		return nil, fmt.Errorf("error decoding response from storage service: %v", err)
	}
	return requestLogs, nil
}

func (s *StorageServiceClient) GetUserDetails(userId string, token string) (*types.GetUserDetailsResponse, error) {
	host := config.Configs.StorageSvcConfig.Host
	port := strconv.Itoa(config.Configs.StorageSvcConfig.Port)
	// connect to storage service
	body, err := json.Marshal(types.GetUserDetailsRequest{
		UserId:         userId,
		RequesterToken: token,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/user-details", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error connecting to storage service: %v", err)
	}
	defer res.Body.Close()
	userDetails := &types.GetUserDetailsResponse{}
	err = json.NewDecoder(res.Body).Decode(userDetails)
	if err != nil {
		return nil, fmt.Errorf("error decoding response from storage service: %v", err)
	}
	return userDetails, nil
}
