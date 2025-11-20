package services

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"cocode/backend/config"
	"cocode/backend/models"
)

var (
	users    = make(map[string]*models.User)
	sessions = make(map[string]*models.Session)
	mu       sync.RWMutex
)

// LoadUsers 从文件加载用户
func LoadUsers() error {
	file, err := os.Open(config.AppConfig.Auth.UsersFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 3 {
			continue
		}

		user := &models.User{
			Username:    parts[0],
			Password:    parts[1],
			DisplayName: parts[2],
		}
		users[user.Username] = user
	}

	return scanner.Err()
}

// Authenticate 验证用户
func Authenticate(username, password string) (*models.User, error) {
	mu.RLock()
	user, exists := users[username]
	mu.RUnlock()

	if !exists {
		return nil, errors.New("用户不存在")
	}

	if user.Password != password {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// CreateSession 创建会话
func CreateSession(username string) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	session := &models.Session{
		SessionID: sessionID,
		Username:  username,
		ExpireAt:  time.Now().Add(time.Hour * time.Duration(config.AppConfig.Auth.SessionTimeout)),
	}

	mu.Lock()
	sessions[sessionID] = session
	mu.Unlock()

	return sessionID, nil
}

// ValidateSession 验证会话
func ValidateSession(sessionID string) (*models.Session, error) {
	mu.RLock()
	session, exists := sessions[sessionID]
	mu.RUnlock()

	if !exists {
		return nil, errors.New("会话不存在")
	}

	if time.Now().After(session.ExpireAt) {
		mu.Lock()
		delete(sessions, sessionID)
		mu.Unlock()
		return nil, errors.New("会话已过期")
	}

	return session, nil
}

// DeleteSession 删除会话
func DeleteSession(sessionID string) {
	mu.Lock()
	delete(sessions, sessionID)
	mu.Unlock()
}

// generateSessionID 生成会话ID
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GetAllUsers 获取所有用户（仅管理员）
func GetAllUsers() []*models.User {
	mu.RLock()
	defer mu.RUnlock()

	userList := make([]*models.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	return userList
}

// SaveUsers 保存用户到文件
func SaveUsers() error {
	mu.RLock()
	defer mu.RUnlock()

	file, err := os.Create(config.AppConfig.Auth.UsersFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, user := range users {
		line := user.Username + ":" + user.Password + ":" + user.DisplayName + "\n"
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}

// CreateUser 创建用户
func CreateUser(username, password, displayName string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[username]; exists {
		return errors.New("用户已存在")
	}

	users[username] = &models.User{
		Username:    username,
		Password:    password,
		DisplayName: displayName,
	}

	return nil
}

// UpdateUser 更新用户
func UpdateUser(username, password, displayName string) error {
	mu.Lock()
	defer mu.Unlock()

	user, exists := users[username]
	if !exists {
		return errors.New("用户不存在")
	}

	if password != "" {
		user.Password = password
	}
	if displayName != "" {
		user.DisplayName = displayName
	}

	return nil
}

// DeleteUser 删除用户
func DeleteUser(username string) error {
	mu.Lock()
	defer mu.Unlock()

	if username == "admin" {
		return errors.New("不能删除管理员账号")
	}

	if _, exists := users[username]; !exists {
		return errors.New("用户不存在")
	}

	delete(users, username)
	return nil
}
