package service

import (
	"authz/config"
	"authz/storage"
	"authz/types"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Verify(t *testing.T) {
	setup()
	defer teardown()
	authService := &AuthenticationService{}
	testUserDetails := &types.UserDetails{
		UserId:    "u_id",
		UserEmail: "testuser",
	}
	ss, err := authService.GenerateToken(testUserDetails)
	assert.NoError(t, err, "GenerateToken() should not throw error")
	assert.NotEmpty(t, ss, "generated token shouldn't be empty")

	ok, err := authService.VerifyToken(ss)
	assert.NoError(t, err)
	assert.True(t, ok)

	config.Configs.AuthSecretKey = "newsecretkey"

	ok, err = authService.VerifyToken(ss)
	assert.Error(t, err)
	assert.False(t, ok)
}

func Test_GenerateToken(t *testing.T) {
	setup()
	defer teardown()
	authService := &AuthenticationService{}
	testUserDetails := &types.UserDetails{
		UserId:    "u_id",
		UserEmail: "testuser",
	}
	ss, err := authService.GenerateToken(testUserDetails)
	assert.NoError(t, err, "GenerateToken() should not throw error")
	assert.NotEmpty(t, ss, "generated token shouldn't be empty")
}

func Test_Login_Success(t *testing.T) {
	setup()
	defer teardown()
	authService := &AuthenticationService{}
	mockedStorage := storage.StorageSvc.(*storage.IStorageServiceMock)
	mockPassword := "mypass"
	mockedStorage.GetUserHashFunc = func(email string) (string, error) {
		h := sha256.New()
		_, err := h.Write([]byte(mockPassword))
		if err != nil {
			return "", fmt.Errorf("error getting hash: %v", err)
		}
		return hex.EncodeToString(h.Sum(nil)), nil
	}
	mockedStorage.GetUserIdFunc = func(s string) (string, error) {
		return "aadsiofasd08", nil
	}

	loginReq := &types.UserLoginRequest{
		UserEmail:    "test@gmail.com",
		UserPassword: mockPassword,
	}

	tkn, err := authService.Login(loginReq)
	assert.NoError(t, err, "should succeed")
	assert.NotEmpty(t, tkn, "should not be empty")
	assert.Equal(t, 1, len(mockedStorage.GetUserHashCalls()))
	assert.Equal(t, 1, len(mockedStorage.GetUserIdCalls()))
}

func Test_Login_Fail(t *testing.T) {
	setup()
	defer teardown()
	authService := &AuthenticationService{}
	mockedStorage := storage.StorageSvc.(*storage.IStorageServiceMock)
	mockPassword := "mypass"
	mockedStorage.GetUserHashFunc = func(email string) (string, error) {
		return "", fmt.Errorf("mock error")
	}
	mockedStorage.GetUserIdFunc = func(s string) (string, error) {
		return "aadsiofasd08", nil
	}

	loginReq := &types.UserLoginRequest{
		UserEmail:    "test@gmail.com",
		UserPassword: mockPassword,
	}

	tkn, err := authService.Login(loginReq)
	assert.Error(t, err, "should not succeed")
	assert.Empty(t, tkn, "should be empty")
	assert.Equal(t, 1, len(mockedStorage.GetUserHashCalls()))
	assert.Equal(t, 0, len(mockedStorage.GetUserIdCalls()))
}

func Test_SignUp_ShouldErrorWhenExistingUserFound(t *testing.T) {
	setup()
	defer teardown()
	authService := &AuthenticationService{}
	storageMock := storage.StorageSvc.(*storage.IStorageServiceMock)
	storageMock.GetUserIdFunc = func(s string) (string, error) {
		return "", nil
	}
	testSignupReq := &types.UserSignupRequest{
		UserEmail:    "test@mail.com",
		UserPassword: "testpass",
	}
	tkn, err := authService.Signup(testSignupReq)
	assert.Error(t, err, "should error when existing email account found")
	assert.Equal(t, "user already exists", err.Error())
	assert.Empty(t, tkn)
	assert.Equal(t, 1, len(storageMock.GetUserIdCalls()))
}

func setup() {
	config.Configs = &config.Config{}
	config.Configs.DBConfig = &config.DBConfig{}
	config.Configs.AuthSecretKey = "test_secret_key"
	storage.StorageSvc = &storage.IStorageServiceMock{}
}

func teardown() {

}
