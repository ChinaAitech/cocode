#!/bin/bash

set -e

echo "========================================"
echo "   协同编程平台 - 构建脚本"
echo "========================================"

# 检查必要工具
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "错误: 未找到 $1，请先安装"
        exit 1
    fi
}

echo "检查依赖工具..."
check_command node
check_command npm
check_command go
check_command g++

echo "✓ 所有依赖工具已安装"
echo ""

# 构建前端
echo "步骤 1/3: 构建前端..."
cd frontend

if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi

echo "构建前端资源..."
npm run build

cd ..
echo "✓ 前端构建完成"
echo ""

# 复制前端资源到后端
echo "步骤 2/3: 嵌入前端资源..."
rm -rf backend/embedded/dist
mkdir -p backend/embedded
cp -r frontend/dist backend/embedded/dist
echo "✓ 前端资源已嵌入"
echo ""

# 构建后端
echo "步骤 3/3: 构建后端..."
cd backend
go mod download
go build -o ../cocode main.go
cd ..
echo "✓ 后端构建完成"
echo ""

# 创建必要的目录
mkdir -p data/temp

echo "========================================"
echo "   构建成功!"
echo "========================================"
echo ""
echo "运行以下命令启动服务器:"
echo "  ./cocode"
echo ""
echo "或指定配置文件:"
echo "  ./cocode config.toml"
echo ""
