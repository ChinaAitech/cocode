package handlers

import (
	"encoding/json"
	"net/http"

	"cocode/backend/services"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	SessionID   string `json:"sessionId,omitempty"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// HandleLogin 处理登录
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "请求格式错误",
		})
		return
	}

	// 验证用户
	user, err := services.Authenticate(req.Username, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 创建会话
	sessionID, err := services.CreateSession(user.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "创建会话失败",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success:     true,
		Message:     "登录成功",
		SessionID:   sessionID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
	})
}

// HandleLogout 处理登出
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.Header.Get("X-Session-ID")
	if sessionID != "" {
		services.DeleteSession(sessionID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "登出成功",
	})
}
