package web

import (
	"log"
	"server/service"
	"server/types"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// authorization - request access
	r.GET("/user-details", authMiddleware, getUserDetails)
	r.GET("/request-logs", authMiddleware, getRequestLogs)

	return r
}

func getRequestLogs(c *gin.Context) {
	userId := c.GetString("user-id")
	if userId == "" {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   "user-id not found in context",
		})
		return
	}
	req := &types.GetRequestLogsRequest{
		RequesterEmail: userId,
	}
	resp, err := service.Svc.GetAccessLogs(*req)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, resp)
}

func getUserDetails(c *gin.Context) {
	userId := c.GetString("user-id")
	userToken := c.GetHeader("Authorization")
	userDetails, err := service.Svc.GetUserDetails(userId, userToken)
	if err != nil {
		log.Printf("error in getting user details: %v\n", err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, userDetails)
}
