package handlers

import (
	"encoding/json"
	"net/http"

	"cocode/backend/services"
)

// HandleGetUsers 获取用户列表（仅管理员）
func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	// 验证会话
	sessionID := r.Header.Get("X-Session-ID")
	session, err := services.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 仅管理员可以查看用户列表
	if session.Username != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	users := services.GetAllUsers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"users":   users,
	})
}

// HandleCreateUser 创建用户（仅管理员）
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// 验证会话
	sessionID := r.Header.Get("X-Session-ID")
	session, err := services.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 仅管理员可以创建用户
	if session.Username != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "用户名和密码不能为空",
		})
		return
	}

	if req.DisplayName == "" {
		req.DisplayName = req.Username
	}

	if err := services.CreateUser(req.Username, req.Password, req.DisplayName); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 保存到文件
	if err := services.SaveUsers(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "保存用户失败: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "用户创建成功",
	})
}

// HandleUpdateUser 更新用户（仅管理员）
func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// 验证会话
	sessionID := r.Header.Get("X-Session-ID")
	session, err := services.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 仅管理员可以更新用户
	if session.Username != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "用户名不能为空",
		})
		return
	}

	if err := services.UpdateUser(req.Username, req.Password, req.DisplayName); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 保存到文件
	if err := services.SaveUsers(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "保存用户失败: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "用户更新成功",
	})
}

// HandleDeleteUser 删除用户（仅管理员）
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// 验证会话
	sessionID := r.Header.Get("X-Session-ID")
	session, err := services.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 仅管理员可以删除用户
	if session.Username != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "用户名不能为空",
		})
		return
	}

	if err := services.DeleteUser(req.Username); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 保存到文件
	if err := services.SaveUsers(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "保存用户失败: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "用户删除成功",
	})
}
