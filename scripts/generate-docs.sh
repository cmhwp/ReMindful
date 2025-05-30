#!/bin/bash

# 生成Swagger API文档
echo "正在生成API文档..."

# 确保swag工具已安装
if ! command -v swag &> /dev/null; then
    echo "swag工具未安装，正在安装..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# 生成文档
swag init -g cmd/server/main.go -d . -o docs

echo "API文档生成完成！"
echo "启动服务器后可访问: http://localhost:8080/swagger/index.html" 