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
