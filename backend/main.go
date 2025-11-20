package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"cocode/backend/config"
	"cocode/backend/handlers"
	"cocode/backend/services"
)

//go:embed embedded/dist
var embeddedFiles embed.FS

func main() {
	// 加载配置
	configPath := "config.toml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 加载用户数据
	if err := services.LoadUsers(); err != nil {
		log.Fatalf("加载用户数据失败: %v", err)
	}

	// 创建临时目录
	if err := os.MkdirAll(config.AppConfig.Compiler.TempDir, 0755); err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}

	// 初始化上传目录
	if err := handlers.InitUploadDir(); err != nil {
		log.Fatalf("创建上传目录失败: %v", err)
	}

	// 初始化WebSocket Hub
	handlers.InitWebSocketHub()

	// 设置路由
	mux := http.NewServeMux()

	// API路由
	mux.HandleFunc("/api/login", handlers.HandleLogin)
	mux.HandleFunc("/api/logout", handlers.HandleLogout)
	mux.HandleFunc("/api/upload", handlers.HandleFileUpload)
	mux.HandleFunc("/api/download/code", handlers.HandleCodeDownload)
	mux.HandleFunc("/uploads/", handlers.HandleFileServe)
	mux.HandleFunc("/ws", handlers.HandleWebSocket)

	// 用户管理API（仅管理员）
	mux.HandleFunc("/api/users", handlers.HandleGetUsers)
	mux.HandleFunc("/api/users/create", handlers.HandleCreateUser)
	mux.HandleFunc("/api/users/update", handlers.HandleUpdateUser)
	mux.HandleFunc("/api/users/delete", handlers.HandleDeleteUser)

	// 静态文件服务
	var staticFS http.FileSystem

	// 尝试使用嵌入的文件系统
	distFS, err := fs.Sub(embeddedFiles, "embedded/dist")
	if err == nil {
		staticFS = http.FS(distFS)
		log.Println("使用嵌入的静态资源")
	} else {
		// 开发模式：使用本地文件
		if _, err := os.Stat(config.AppConfig.Server.StaticPath); err == nil {
			staticFS = http.Dir(config.AppConfig.Server.StaticPath)
			log.Printf("使用本地静态资源: %s", config.AppConfig.Server.StaticPath)
		} else {
			log.Println("警告: 未找到静态资源，前端将无法访问")
		}
	}

	if staticFS != nil {
		mux.Handle("/", http.FileServer(staticFS))
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	log.Printf("服务器启动: http://%s", addr)
	log.Printf("WebSocket地址: ws://%s/ws", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
