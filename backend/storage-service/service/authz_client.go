package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"storage-service/config"
	"storage-service/types"
)

var AuthzClientSvc IAuthzClientService

type IAuthzClientService interface {
	VerifyAccessRequest(requesterToken string, userId string, resources []types.ResourceType) (bool, string, error)
}

func SetAuthzClientService(svc IAuthzClientService) {
	AuthzClientSvc = svc
}

func init() {
	AuthzClientSvc = &AuthzClientService{}
}

type AuthzClientService struct {
}

func (as *AuthzClientService) VerifyAccessRequest(requesterToken string, userId string, resources []types.ResourceType) (bool, string, error) {
	resp, err := http.Get(config.Configs.AuthzConfig.Host + ":" + config.Configs.AuthzConfig.Port + "/request-access")
	if err != nil {
		return false, "", fmt.Errorf("error in making request to authz service: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("error in reading response from authz service: %v", err)
	}
	accessReqResp := &types.RequestAccessResponse{}
	err = json.Unmarshal(body, accessReqResp)
	if err != nil {
		return false, "", fmt.Errorf("error in unmarshalling response from authz service: %v", err)
	}
	if accessReqResp.AccessRequest {
		return true, accessReqResp.Message, nil
	} else {
		return false, accessReqResp.Message, nil
	}
}