#!/bin/bash

# 更新包列表并安装musl-dev和gcc
apt-get update
apt-get install -y musl-dev gcc

# 设置GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct

# 启用CGO
export CGO_ENABLED=1

# 格式化Go代码并整理依赖
go fmt
go mod tidy

# 编译Go代码
go build -o wechat-bot .
