# 快速启动指南

## 一键构建运行

```bash
# 1. 构建项目（首次运行需要）
./build.sh

# 2. 启动服务器
./cocode

# 3. 打开浏览器访问
http://localhost:8080
```

## 默认登录账号

| 用户名 | 密码 | 显示名称 |
|--------|------|----------|
| admin  | admin123 | 管理员 |
| userA  | passwordA | 用户A |
| userB  | passwordB | 用户B |
| userC  | passwordC | 用户C |

## 测试协同功能

1. **测试多人编辑**
   - 用不同浏览器（或隐私窗口）登录不同账号
   - 在代码编辑器中输入代码
   - 观察其他窗口的实时同步

2. **测试编译运行**
   ```cpp
   #include <iostream>
   using namespace std;

   int main() {
       int a, b;
       cin >> a >> b;
       cout << "Sum: " << (a + b) << endl;
       return 0;
   }
   ```
   - 在输入区填写: `5 10`
   - 点击"编译并运行"
   - 查看输出区和日志区

3. **测试聊天功能**
   - 在右侧聊天面板发送消息
   - 所有在线用户都能收到

## 常见问题

### Q: 构建失败提示找不到 g++

```bash
# Ubuntu/Debian
sudo apt-get install g++

# CentOS/RHEL
sudo yum install gcc-c++

# macOS
xcode-select --install
```

### Q: 前端构建超时

网络问题导致 npm 安装慢，可以使用国内镜像：

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
npm install
```

### Q: 端口8080被占用

修改 `config.toml` 中的端口：

```toml
[server]
port = 8888  # 改为其他端口
```

## 开发模式

如果需要修改代码并快速测试：

```bash
# 终端1: 启动后端（开发模式）
cd backend
go run main.go ../config.toml

# 终端2: 启动前端（开发模式）
cd frontend
npm run dev

# 访问 http://localhost:3000
```

开发模式下，前端修改会自动热重载。

## 部署到服务器

```bash
# 1. 在本地构建
./build.sh

# 2. 上传文件到服务器
scp cocode config.toml user@server:/path/to/app/
scp -r data user@server:/path/to/app/

# 3. 在服务器上运行
ssh user@server
cd /path/to/app
./cocode
```

### 后台运行（推荐使用 systemd）

创建 `/etc/systemd/system/cocode.service`:

```ini
[Unit]
Description=C++ Collaboration Platform
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/app
ExecStart=/path/to/app/cocode
Restart=always

[Install]
WantedBy=multi-user.target
```

然后启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable cocode
sudo systemctl start cocode
sudo systemctl status cocode
```

## 性能优化

### 限制资源使用

修改 `config.toml`:

```toml
[compiler]
compile_timeout = 10  # 编译超时（秒）
run_timeout = 5       # 运行超时（秒）
```

### 清理临时文件

临时文件保存在 `data/temp/` 目录，定期清理：

```bash
# 添加到 crontab
0 2 * * * find /path/to/app/data/temp -type f -mtime +1 -delete
```

## 安全建议

⚠️ **生产环境必读**

1. **修改默认密码**: 编辑 `data/users.txt`
2. **使用 HTTPS**: 配置 nginx 反向代理
3. **限制访问**: 配置防火墙规则
4. **沙箱隔离**: 使用 Docker 容器运行
5. **资源限制**: 使用 cgroup 限制 CPU/内存

### Nginx 反向代理示例

```nginx
server {
    listen 443 ssl;
    server_name cocode.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 技术支持

遇到问题？

1. 查看日志输出
2. 检查 `data/temp/` 目录权限
3. 确认 g++ 编译器可用
4. 查看 README.md 详细文档
