package web

import (
	"encoding/json"
	"io"
	"storage-service/service"
	"storage-service/types"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// authG := r.Group("/auth/")
	// authG.GET("/verify", verify)
	// authG.POST("/gen-token", generateToken)
	// authG.POST("/login", login)
	// authG.POST("/signup", signup)

	// authorization - request access
	// r.GET("/request-access", requestAccess)
	r.GET("/data", getData)
	r.PUT("/data", writeData)

	return r
}

func writeData(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	writeDataReq := &types.WriteDataRequest{}
	err = json.Unmarshal(bytes, writeDataReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   err.Error(),
		})
		return
	}
	ok := false
	for _, resource := range types.ResourceList {
		if _, ok1 := writeDataReq.Data[resource]; ok1 {
			ok = true
		}
	}
	if !ok {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   "Invalid resource type",
		})
		return
	}

	resp, err := service.Svc.WriteData(writeDataReq)
	if err != nil {
		if err == types.ErrAccessDenied {
			c.JSON(403, gin.H{
				"error": "Access Denied",
				"msg":   resp.Msg,
			})
			return
		}
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, resp)
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
	ok := true
	for _, resource := range getDataReq.Resources {
		if _, ok1 := types.ResourceIdToString[resource]; !ok1 {
			ok = false
			break
		}
	}
	if !ok {
		c.JSON(400, gin.H{
			"error": "Bad Request",
			"msg":   "Invalid resource type",
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
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, resp)
}
