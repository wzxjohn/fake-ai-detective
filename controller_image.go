package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func handleFakeImage(c *gin.Context) {
	traceID := c.Param("id")

	var req *TrackedRequest
	var exists bool
	mutex.Lock()
	if req, exists = recordedIPs[traceID]; !exists {
		mutex.Unlock()
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	mutex.Unlock()

	// 处理用户代理信息
	var message string
	userAgent := c.GetHeader("User-Agent")
	if strings.Contains(userAgent, "IPS") {
		message = "可能是 Azure"
	} else if strings.Contains(userAgent, "OpenAI") {
		message = "可能是 OpenAI"
	} else if strings.Contains(userAgent, "Go-http-client") {
		message = "可能是 Go 中转"
	} else {
		message = "未知，可能来自逆向"
	}

	recordMessage(traceID, time.Now().Unix(), c.ClientIP(), c.GetHeader("User-Agent"), message, c.Request.Header)

	// 返回 WebP 图片
	c.Data(http.StatusOK, "image/webp", req.Image)
}