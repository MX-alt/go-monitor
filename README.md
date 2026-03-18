# Go System Monitor (Distributed Architecture)

这是一个基于 Go 语言开发的分布式系统监控原型，实现了从底层数据采集到容器化编排的全栈架构。本项目模拟了生产环境中的高可用部署方案。

## 🏗️ 架构概览

该项目采用“双容器”编排模式：
1. **Frontend Proxy (Nginx)**: 充当负载均衡与静态资源网关，监听 `8080` 端口。
2. **Backend Service (Go)**: 负责系统资源（如磁盘、内存）的实时监控数据采集，监听 `8081` 端口。



## 🚀 快速开始 (Getting Started)

本项目已完全容器化，支持一键部署。

### 前置要求
- Docker Engine >= 20.10.0
- Docker Compose V2 (推荐使用 `docker compose` 命令)

### 部署命令
```bash
# 克隆仓库
git clone [https://github.com/MX-alt/go-system-monitor.git](https://github.com/MX-alt/go-system-monitor.git)
cd go-system-monitor

# 启动全栈服务
docker compose up -d