package main

import (
	"encoding/base64"
	"fake-ai-detective/config"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleIndex(c *gin.Context) {
	apiPrefix := config.GetConfig().APIPrefix
	data := gin.H{}
	if apiPrefix != "/" {
		data["APIPrefix"] = template.HTML(apiPrefix)
	} else {
		data["APIPrefix"] = ""
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func handleStart(c *gin.Context) {
	startReq := &StartRequest{}
	err := c.BindJSON(startReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := strings.TrimSpace(startReq.URL)
	key := strings.TrimSpace(startReq.Key)
	model := strings.TrimSpace(startReq.Model)
	if url == "" || key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	if model == "" {
		model = "gpt-4o"
	}

	traceID := uuid.NewString()

	req := &TrackedRequest{
		Timestamp: time.Now(),
		Image:     nil,
		Results:   make([]*TrackedResult, 0, 8),
	}

	imgBytes, err := generateCaptchaImage(traceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Image = imgBytes

	mutex.Lock()
	recordedIPs[traceID] = req
	mutex.Unlock()

	// 异步发送 POST 请求
	go sendOpenAIRequest(url, key, model, traceID)

	startInfo := &StartInfo{
		ID:    traceID,
		Image: "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBytes),
	}

	c.JSON(http.StatusOK, &StartResponse{Data: startInfo})
}

func handleResult(c *gin.Context) {
	traceID := c.Param("id")

	mutex.RLock()
	req, exists := recordedIPs[traceID]
	mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusOK, []string{})
		return
	}

	mutex.RLock()
	finished := req.Finished
	results := req.Results
	mutex.RUnlock()

	c.JSON(http.StatusOK, &ResultResponse{Data: &DetectResult{
		Finished: finished,
		Results:  results,
	}})
}

func handleResponse(c *gin.Context) {
	traceID := c.Param("id")

	mutex.RLock()
	req, exists := recordedIPs[traceID]
	mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusOK, []string{})
		return
	}

	mutex.RLock()
	response := req.Response
	mutex.RUnlock()

	c.JSON(http.StatusOK, &TargetResponse{Data: response})
}
