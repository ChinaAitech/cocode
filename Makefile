.PHONY: all frontend backend clean run dev

# 默认目标：构建前端和后端
all: frontend backend

# 构建前端
frontend:
	@echo "正在构建前端..."
	cd frontend && npm install && npm run build
	@echo "复制前端资源到后端嵌入目录..."
	mkdir -p backend/embedded
	rm -rf backend/embedded/dist
	cp -r frontend/dist backend/embedded/dist
	@echo "前端构建完成!"

# 构建后端
backend:
	@echo "正在构建后端..."
	cd backend && go mod download && go build -o ../cocode main.go
	@echo "后端构建完成!"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf backend/embedded/dist
	rm -f cocode
	rm -rf data/temp/*
	@echo "清理完成!"

# 运行应用（开发模式）
dev:
	@echo "启动开发模式..."
	cd backend && go run main.go ../config.toml

# 运行已构建的应用
run:
	@echo "启动应用..."
	./cocode config.toml

# 安装依赖
deps:
	@echo "安装依赖..."
	cd frontend && npm install
	cd backend && go mod download
	@echo "依赖安装完成!"

# 快速构建（仅后端，用于开发）
quick:
	@echo "快速构建后端..."
	cd backend && go build -o ../cocode main.go
	@echo "快速构建完成!"

# 构建 Linux 版本
linux:
	@echo "正在构建 Linux 版本..."
	cd backend && GOOS=linux GOARCH=amd64 go build -o ../cocode-linux main.go
	@echo "Linux 版本构建完成! (cocode-linux)"

# 构建 Windows 版本
win:
	@echo "正在构建 Windows 版本..."
	cd backend && GOOS=windows GOARCH=amd64 go build -o ../cocode-win.exe main.go
	@echo "Windows 版本构建完成! (cocode-win.exe)"
