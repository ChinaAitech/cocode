package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cocode/backend/services"
)

const (
	maxUploadSize = 50 * 1024 * 1024 // 50MB
	uploadDir     = "./data/uploads"
)

// InitUploadDir 初始化上传目录
func InitUploadDir() error {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("创建上传目录失败: %v", err)
	}
	return nil
}

// HandleFileUpload 处理文件上传
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	// 验证会话
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		sessionID = r.URL.Query().Get("session")
	}

	_, err := services.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 限制上传大小
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "文件太大", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "获取文件失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 生成唯一文件名
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), sanitizeFilename(strings.TrimSuffix(handler.Filename, ext)), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("创建文件失败: %v", err)
		http.Error(w, "保存文件失败", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("写入文件失败: %v", err)
		http.Error(w, "保存文件失败", http.StatusInternalServerError)
		return
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success":true,"fileUrl":"%s","fileName":"%s","fileSize":%d}`,
		fileURL, handler.Filename, handler.Size)
}

// HandleFileServe 提供文件下载服务
func HandleFileServe(w http.ResponseWriter, r *http.Request) {
	// 获取文件名
	filename := strings.TrimPrefix(r.URL.Path, "/uploads/")
	if filename == "" {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}

	// 防止路径遍历攻击
	filename = filepath.Base(filename)
	filePath := filepath.Join(uploadDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}

	// 提供文件
	http.ServeFile(w, r, filePath)
}

// sanitizeFilename 清理文件名
func sanitizeFilename(filename string) string {
	// 移除特殊字符
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		"..", "_",
		" ", "_",
	)
	return replacer.Replace(filename)
}
