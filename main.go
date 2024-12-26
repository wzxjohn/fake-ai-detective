package main

import (
	"embed"
	"fake-ai-detective/config"
	"html/template"
	"log"
	"path"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 全局变量
var (
	//go:embed web/*
	templateFS embed.FS

	recordedIPs = make(map[string]*TrackedRequest)
	mutex       sync.RWMutex
	baseURL     string
)

func main() {
	cfg := config.GetConfig()

	r := gin.Default()
	err := r.SetTrustedProxies(cfg.TrustedProxies)
	if err != nil {
		panic(err)
	}

	// 从 embed.FS 加载模板
	tmpl := template.Must(template.ParseFS(templateFS, "web/*.html"))
	r.SetHTMLTemplate(tmpl)
	r.GET("/", handleIndex)

	// CORS 中间件配置
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowCredentials = true
	corsCfg.AllowMethods = []string{"*"}
	corsCfg.AllowHeaders = []string{"*"}
	r.Use(cors.New(corsCfg))

	baseURL = "https://" + cfg.Domain + path.Clean("/"+cfg.ImagePrefix) + "/"
	log.Println("Base URL: ", baseURL)

	// 路由设置
	apiGroup := r.Group(cfg.APIPrefix)
	apiGroup.POST("/start", handleStart)
	apiGroup.GET("/result/:id", handleResult)

	imgGroup := r.Group(cfg.ImagePrefix)
	imgGroup.GET("/:id", handleFakeImage)

	// 启动清理旧数据的 goroutine
	go cleanupOldIPs()

	// 启动服务器
	r.Run()
}
