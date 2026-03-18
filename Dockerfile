# 使用轻量级的 Alpine Linux 作为基础，它只有 5MB 左右，非常符合 SRE 的极简主义
FROM golang:1.26-alpine

# 设置容器内部的工作目录
WORKDIR /app

# 将当前目录下的 go.mod 复制进去（先复制这个可以利用 Docker 的缓存机制，加速后续构建）
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制剩下的所有源代码
COPY . .

# 编译成名为 "monitor" 的二进制文件
RUN go build -o monitor main.go

# 声明容器运行时监听 8081 端口
EXPOSE 8081

# 启动命令
CMD ["./monitor"]