package main

import (
	"net/http"
	"time"
)

func recordMessage(traceID string, ts int64, ip, ua, message string, header http.Header) {
	mutex.Lock()
	defer mutex.Unlock()
	result := &TrackedResult{
		Timestamp: ts,
		IP:        ip,
		UserAgent: ua,
		Message:   message,
		Header:    header,
	}
	if req, exists := recordedIPs[traceID]; exists {
		req.Results = append(req.Results, result)
	}
}

func markFinished(traceID string) {
	mutex.Lock()
	defer mutex.Unlock()
	if req, exists := recordedIPs[traceID]; exists {
		req.Finished = true
	}
}

func cleanupOldIPs() {
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {
		mutex.Lock()
		now := time.Now()
		for id, req := range recordedIPs {
			if now.Sub(req.Timestamp) > 3*time.Minute {
				delete(recordedIPs, id)
			}
		}
		mutex.Unlock()
	}
}
