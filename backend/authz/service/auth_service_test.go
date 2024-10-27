package service

import (
	"authz/config"
	"authz/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Verify(t *testing.T) {
	setup()
	defer restore()
	authService := &AuthenticationService{}
	testUserDetails := &types.UserDetails{
		UserId:   "u_id",
		Username: "testuser",
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
	defer restore()
	authService := &AuthenticationService{}
	testUserDetails := &types.UserDetails{
		UserId:   "u_id",
		Username: "testuser",
	}
	ss, err := authService.GenerateToken(testUserDetails)
	assert.NoError(t, err, "GenerateToken() should not throw error")
	assert.NotEmpty(t, ss, "generated token shouldn't be empty")
}



func setup() {
	config.Configs = &config.Config{}
	config.Configs.DBConfig = &config.DBConfig{}
	config.Configs.AuthSecretKey = "test_secret_key"
}

func restore() {

}
