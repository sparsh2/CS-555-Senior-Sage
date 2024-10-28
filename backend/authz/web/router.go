package web

import (
	"authz/service"
	"authz/types"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authG := r.Group("/auth/")
	authG.PUT("/verify", verify)
	authG.POST("/gen-token", generateToken)
	authG.POST("/login", login)

	return r
}

func login(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	loginReq := &types.UserLoginRequest{}
	err = json.Unmarshal(bytes, loginReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}

	token, err := service.AuthService.Login(loginReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
		"msg": "login successful",
	})
}

func generateToken(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	tokenGenReq := &types.UserDetails{}
	err = json.Unmarshal(bytes, tokenGenReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	tokenString, err := service.AuthService.GenerateToken(tokenGenReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error generating token",
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

func verify(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	authReq := &types.AuthVerifyRequest{}
	err = json.Unmarshal(bytes, authReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}

	ok, err := service.AuthService.VerifyToken(authReq.JWTToken)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Unknown internal error",
			"msg": err.Error(),
		})
		return
	}

	if ok {
		c.JSON(200, gin.H{
			"valid": "true",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"valid": "false",
		})
	}
}
