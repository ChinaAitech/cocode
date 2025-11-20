# C++ 协同编程平台 - 功能清单

## 已实现功能 ✅

### 1. 聊天系统增强
- ✅ 不同用户显示不同颜色（20种颜色自动分配）
- ✅ 消息布局优化（自己右侧，他人左侧）
- ✅ 表情选择器（56个常用表情）
- ✅ 图片上传分享
- ✅ 文件上传分享（最大50MB）
- ✅ 支持 Shift+Enter 换行

### 2. 代码管理
- ✅ 实时协同编辑
- ✅ 代码下载功能（.cpp文件）
- ✅ Monaco编辑器集成

### 3. 编译运行
- ✅ C++17 编译支持
- ✅ 标准输入支持
- ✅ 编译日志显示
- ✅ 运行结果输出

### 4. 用户系统
- ✅ 简单认证系统
- ✅ 会话管理
- ✅ 在线用户列表

## 待实现功能 📋

### 高优先级
- [ ] 输入输出同步（所有用户看到相同的输入/输出/日志）
- [ ] 编译记录（显示谁执行了编译）
- [ ] 输入区文件上传（上传文件内容作为输入）
- [ ] 标准答案区（用于OJ判题）
- [ ] 输出对比功能（逐行对比）

### 中优先级
- [ ] 消息引用回复
- [ ] 100个用户账号支持
- [ ] Admin用户管理界面
- [ ] 用户权限控制

### 低优先级
- [ ] 代码历史版本
- [ ] 代码模板库
- [ ] 多语言支持（不仅C++）
- [ ] 性能优化

## 技术栈

### 后端
- Go 1.21
- Gorilla WebSocket
- 文件存储（用户/会话）

### 前端
- Vue 3
- Element Plus
- Monaco Editor
- Vite

## API 接口

### 已实现
- `POST /api/login` - 用户登录
- `POST /api/logout` - 用户登出
- `POST /api/upload` - 文件上传
- `GET /api/download/code` - 代码下载
- `GET /uploads/:filename` - 文件下载
- `WS /ws` - WebSocket连接

### 计划中
- `GET /api/users` - 获取用户列表（Admin）
- `POST /api/users` - 创建用户（Admin）
- `DELETE /api/users/:id` - 删除用户（Admin）
- `GET /api/compile/history` - 编译历史

## 文件结构

```
cocode/
├── backend/
│   ├── config/          # 配置管理
│   ├── handlers/        # HTTP/WebSocket处理器
│   │   ├── auth.go
│   │   ├── websocket.go
│   │   ├── upload.go   # 文件上传
│   │   └── code.go     # 代码下载
│   ├── models/          # 数据模型
│   ├── services/        # 业务逻辑
│   └── main.go
├── frontend/
│   └── src/
│       ├── components/
│       │   ├── ChatPanel.vue    # 聊天组件
│       │   └── CodeEditor.vue   # 编辑器组件
│       └── App.vue
├── data/
│   ├── users.txt        # 用户数据
│   ├── uploads/         # 上传文件
│   └── temp/            # 临时编译文件
└── config.toml          # 主配置
```

## 数据存储

### 用户数据 (data/users.txt)
```
username:password:displayname
admin:admin123:管理员
userA:passwordA:用户A
```

### 上传文件 (data/uploads/)
- 自动生成唯一文件名
- 防止路径遍历攻击
- 50MB大小限制

## 安全注意事项

⚠️ **当前版本为演示项目，生产环境需要改进：**
- 密码采用明文存储（需要bcrypt哈希）
- 代码执行无沙箱隔离（建议Docker容器）
- WebSocket无速率限制
- 未实现HTTPS
- 需要加强输入验证

## 更新日志

### 2024-11-20
- 添加聊天多媒体功能（表情、图片、文件）
- 实现代码下载功能
- 优化聊天UI（颜色区分、左右布局）
- 添加文件上传系统（50MB限制）

### 2024-11-18
- 初始版本发布
- 实现基础协同编辑功能
- 实现C++编译运行
- 实现简单聊天系统
