package handlers

import (
	"fmt"
	"net/http"
	"time"

	"cocode/backend/services"
)

// HandleCodeDownload 处理代码下载
func HandleCodeDownload(w http.ResponseWriter, r *http.Request) {
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

	// 获取当前代码
	if hub == nil {
		http.Error(w, "服务未就绪", http.StatusServiceUnavailable)
		return
	}

	codeState := hub.GetCodeState()
	code := codeState.Code
	if code == "" {
		code = "// 空代码\n"
	}

	// 设置下载响应头
	filename := fmt.Sprintf("code_%s.cpp", time.Now().Format("20060102_150405"))
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(code)))

	// 写入代码内容
	w.Write([]byte(code))
}
