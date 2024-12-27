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
	ip := c.ClientIP()
	if IsFromOpenAI(ip) {
		message += "IP 来自 OpenAI, "
	} else if IsFromCloudflare(ip) {
		message += "IP 来自 Cloudflare, "
	} else {
		message += "IP 未知, "
	}

	userAgent := c.GetHeader("User-Agent")
	if strings.Contains(userAgent, "IPS") {
		message += "可能是 Azure"
	} else if strings.Contains(userAgent, "OpenAI") {
		message += "可能是 OpenAI"
	} else if strings.Contains(userAgent, "Go-http-client") {
		message += "可能是 Go 中转"
	} else {
		message += "未知，可能来自逆向"
	}

	recordMessage(traceID, time.Now().Unix(), ip, c.GetHeader("User-Agent"), message, c.Request.Header)

	// 返回图片
	c.Data(http.StatusOK, "image/png", req.Image)
}
