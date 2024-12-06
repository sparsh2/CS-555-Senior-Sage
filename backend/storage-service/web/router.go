package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"storage-service/service"
	"storage-service/types"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.PUT("/preferences", writePreferences)
	r.PUT("/reponses", writeResponses)
	r.PUT("/reminders", writeReminders)
	r.PUT("/question-counter", writeQuestionCounter)
	r.PUT("/chat-history", writeChatHistory)
	r.POST("/data", getData)
	r.POST("/access-logs", getRequestLogs)

	return r
}

func getRequestLogs(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	req := &types.GetRequestLogsRequest{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	logs, err := service.Svc.GetAccessLogs(req.RequesterEmail)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	resp := &types.GetRequestLogsResponse{
		Logs: logs,
	}
	c.JSON(200, resp)
}

func writeChatHistory(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	req := &types.WriteChatHistoryRequest{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	err = service.Svc.WriteChatHistory(req)
	if err == types.ErrAccessDenied {
		c.JSON(403, gin.H{
			"error": "Access Denied",
			"msg":   "Write access to chat history denied",
		})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func writeQuestionCounter(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	req := &types.WriteQuestionCounterRequest{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	err = service.Svc.WriteQuestionCounter(req)
	if err == types.ErrAccessDenied {
		c.JSON(403, gin.H{
			"error": "Access Denied",
			"msg":   "Write access to question counter denied",
		})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func writeReminders(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	req := &types.WriteRemindersRequest{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	if req.Reminders == nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   "Reminders cannot be nil",
		})
		return
	}
	err = service.Svc.WriteReminders(req)
	if err == types.ErrAccessDenied {
		c.JSON(403, gin.H{
			"error": "Access Denied",
			"msg":   "Write access to reminders denied",
		})
		return
	}
	if err != nil {
		log.Printf("error in writing reminders: %v", err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func writeResponses(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	req := &types.WriteResponsesRequest{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	err = service.Svc.WriteResponses(req)
	if err == types.ErrAccessDenied {
		c.JSON(403, gin.H{
			"error": "Access Denied",
			"msg":   "Write access to responses denied",
		})
		return
	}
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func writePreferences(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	writePreferencesReq := &types.WritePreferencesRequest{}
	err = json.Unmarshal(bytes, writePreferencesReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	err = service.Svc.WritePreferences(writePreferencesReq)
	if err == types.ErrAccessDenied {
		c.JSON(403, gin.H{
			"error": "Access Denied",
			"msg":   "Write access to preferences denied",
		})
		return
	}
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func getData(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	getDataReq := &types.GetDataRequest{}
	err = json.Unmarshal(bytes, getDataReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	if getDataReq.RequesterToken == "" {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   "Requester token cannot be nil",
		})
		return
	}

	resp, err := service.Svc.GetData(getDataReq)
	if err != nil {
		if err == types.ErrAccessDenied {
			c.JSON(403, gin.H{
				"error": "Access Denied",
				"msg":   resp.Msg,
			})
			return
		}
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, resp)
}
