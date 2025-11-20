package models

import (
	"time"
)

// User 用户结构
type User struct {
	Username    string `json:"username"`
	Password    string `json:"-"` // 不序列化密码
	DisplayName string `json:"displayName"`
}

// Session 会话结构
type Session struct {
	SessionID string
	Username  string
	ExpireAt  time.Time
}

// WebSocketMessage WebSocket消息结构
type WebSocketMessage struct {
	Type        string      `json:"type"`               // message类型: edit, cursor, chat, compile, run
	Username    string      `json:"username"`           // 发送者用户名
	DisplayName string      `json:"displayName"`        // 发送者显示名称
	Timestamp   int64       `json:"timestamp"`          // 时间戳
	Data        interface{} `json:"data"`               // 消息数据
}

// EditOperation 编辑操作
type EditOperation struct {
	Position int    `json:"position"` // 编辑位置
	Content  string `json:"content"`  // 编辑内容
	Action   string `json:"action"`   // insert, delete, replace
	Length   int    `json:"length"`   // 影响长度
}

// CursorPosition 光标位置
type CursorPosition struct {
	Line   int    `json:"line"`   // 行号
	Column int    `json:"column"` // 列号
	Color  string `json:"color"`  // 光标颜色（区分不同用户）
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Message     string `json:"message"`
	MessageType string `json:"messageType"` // text, emoji, image, file
	FileURL     string `json:"fileUrl,omitempty"`
	FileName    string `json:"fileName,omitempty"`
	FileSize    int64  `json:"fileSize,omitempty"`
	ReplyTo     int    `json:"replyTo,omitempty"` // 引用的消息ID
}

// CompileRequest 编译请求
type CompileRequest struct {
	Code  string `json:"code"`  // C++代码
	Input string `json:"input"` // 输入数据
}

// CompileResult 编译结果
type CompileResult struct {
	Success bool   `json:"success"` // 编译是否成功
	Message string `json:"message"` // 编译信息/错误
	Output  string `json:"output"`  // 运行输出
}

// CodeState 代码状态（用于协同编辑）
type CodeState struct {
	Code    string    `json:"code"`    // 当前代码
	Version int       `json:"version"` // 版本号
	Updated time.Time `json:"updated"` // 更新时间
}

// SharedState 共享状态（输入、输出、日志）
type SharedState struct {
	InputData   string    `json:"inputData"`   // 输入数据
	OutputData  string    `json:"outputData"`  // 输出数据
	CompileLog  string    `json:"compileLog"`  // 编译日志
	Answer      string    `json:"answer"`      // 标准答案
	Updated     time.Time `json:"updated"`     // 更新时间
}

// CompileRecord 编译记录
type CompileRecord struct {
	Username  string    `json:"username"`  // 执行编译的用户
	Timestamp time.Time `json:"timestamp"` // 编译时间
	Success   bool      `json:"success"`   // 是否成功
}
