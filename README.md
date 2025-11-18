# C++ 协同编程平台

一个基于Web的多人协同C++编程平台，支持实时代码同步、在线编译运行和聊天功能。

## 功能特性

- **实时协同编辑**: 多人同时编辑代码，实时同步
- **光标同步**: 查看其他用户的光标位置
- **在线编译运行**: 支持C++17标准，实时编译并运行代码
- **输入输出**: 提供标准输入区和输出区
- **聊天功能**: 内置聊天室，方便团队沟通
- **用户管理**: 简单的文本文件认证系统
- **单文件部署**: 前端嵌入到Go二进制文件中，只需一个bin文件和一个配置文件

## 技术栈

### 前端
- Vue 3
- Element Plus (UI组件库)
- Monaco Editor (代码编辑器)
- WebSocket (实时通信)

### 后端
- Golang
- Gorilla WebSocket
- TOML配置文件

## 系统要求

- Node.js 16+ 和 npm (构建时需要)
- Go 1.21+ (构建时需要)
- g++ (运行时需要，用于编译C++代码)
- Linux 系统 (推荐)

## 快速开始

### 方式一：使用构建脚本

```bash
# 1. 克隆项目
git clone <repository-url>
cd cocode

# 2. 运行构建脚本
./build.sh

# 3. 启动服务器
./cocode
```

### 方式二：使用Makefile

```bash
# 构建整个项目
make all

# 运行应用
make run

# 开发模式（不嵌入前端资源）
make dev
```

### 方式三：手动构建

```bash
# 1. 构建前端
cd frontend
npm install
npm run build
cd ..

# 2. 复制前端资源
mkdir -p backend/embedded
cp -r frontend/dist backend/embedded/dist

# 3. 构建后端
cd backend
go build -o ../cocode main.go
cd ..

# 4. 启动
./cocode
```

## 配置

编辑 `config.toml` 文件来修改配置：

```toml
[server]
host = "0.0.0.0"
port = 8080

[compiler]
compiler = "g++"
compile_flags = ["-std=c++17", "-Wall"]
compile_timeout = 30
run_timeout = 10

[auth]
users_file = "./data/users.txt"
session_timeout = 24

[websocket]
path = "/ws"
ping_interval = 30
```

## 用户管理

用户信息存储在 `data/users.txt` 文件中，格式为：

```
username:password:displayname
```

默认用户：
- admin:admin123:管理员
- userA:passwordA:用户A
- userB:passwordB:用户B
- userC:passwordC:用户C

## 使用说明

1. **登录**: 使用上述用户名和密码登录系统

2. **编辑代码**: 在代码编辑器中编写C++代码，代码会实时同步给所有在线用户

3. **编译运行**:
   - 在"输入区"填写程序需要的标准输入
   - 点击"编译并运行"按钮
   - 编译信息显示在"编译日志"区
   - 运行结果显示在"输出区"

4. **聊天**: 使用右侧聊天面板与其他用户交流

5. **查看在线用户**: 顶部显示当前所有在线用户

## 项目结构

```
cocode/
├── frontend/              # Vue前端项目
│   ├── src/
│   │   ├── components/   # Vue组件
│   │   ├── App.vue       # 主应用
│   │   └── main.js       # 入口文件
│   ├── package.json
│   └── vite.config.js
│
├── backend/              # Go后端
│   ├── config/           # 配置管理
│   ├── handlers/         # HTTP处理器
│   ├── models/           # 数据模型
│   ├── services/         # 业务逻辑
│   ├── embedded/         # 嵌入的前端资源
│   └── main.go           # 入口文件
│
├── data/                 # 数据文件
│   ├── users.txt         # 用户数据
│   └── temp/             # 临时文件（编译）
│
├── config.toml           # 配置文件
├── build.sh              # 构建脚本
├── Makefile              # Make构建文件
└── README.md             # 本文件
```

## 开发

### 前端开发

```bash
cd frontend
npm install
npm run dev  # 开发服务器运行在 http://localhost:3000
```

### 后端开发

```bash
cd backend
go run main.go ../config.toml
```

前端开发服务器会自动代理API和WebSocket请求到后端。

## 清理

```bash
# 清理所有构建文件
make clean

# 或手动清理
rm -rf frontend/dist frontend/node_modules
rm -rf backend/embedded/dist
rm -f cocode
```

## 安全注意事项

⚠️ **警告**: 这是一个演示项目，不建议直接用于生产环境：

1. 密码使用明文存储，生产环境应使用哈希加密
2. 代码执行没有沙箱隔离，存在安全风险
3. WebSocket没有速率限制
4. 没有输入验证和SQL注入防护

如需用于生产环境，请考虑：
- 使用数据库存储用户信息
- 实现密码哈希（bcrypt）
- 添加代码执行沙箱（Docker容器）
- 实现速率限制和CSRF防护
- 添加HTTPS支持

## 常见问题

### 1. 编译失败

确保系统已安装 g++ 编译器：
```bash
sudo apt-get install g++  # Ubuntu/Debian
sudo yum install gcc-c++   # CentOS/RHEL
```

### 2. WebSocket连接失败

检查防火墙设置，确保端口8080开放：
```bash
sudo ufw allow 8080  # Ubuntu
```

### 3. 前端资源404

确保已正确构建并嵌入前端资源：
```bash
make frontend
make backend
```

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 作者

协同编程平台开发团队
