# 🚀 K8s-Go System Monitor Dashboard

[![Kubernetes](https://img.shields.io/badge/Kubernetes-v1.28+-blue.svg)](https://kubernetes.io/)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)

这是一个基于云原生架构的系统实时监控项目。它展示了如何通过 Kubernetes 编排多服务组件，并实现从底层硬件指标采集到前端 UI 自动刷新的完整闭环。

---

## 🏗 核心架构 (Architecture)

该项目采用微服务分层设计，通过 K8s 进行流量治理：

* **Backend (Go)**: 实时计算容器内的磁盘空间占用，通过 `net/http` 暴露 `/api/disk` JSON 接口。
* **Frontend (Nginx + HTML/JS)**: 
    * 作为静态资源服务器分发监控面板。
    * 作为 **Reverse Proxy (反向代理)**，处理前端跨域问题并路由流量至后端服务。
* **Infrastructure (K8s)**:
    * 使用 `Service` (ClusterIP) 实现内部稳定的服务发现。
    * 使用 `ConfigMap` 挂载 Nginx 配置，实现逻辑与配置的分离。

```text
                    +------------------------------------------+
                    |           User Browser (Client)          |
                    |   (Accesses Cloud Shell Port 8080)       |
                    +------------------+-----------------------+
                                       |
                                       | HTTP Request
                                       v
                    +------------------------------------------+
                    |          K8s Service (nginx-service)     |
                    |          (ClusterIP: 8080)               |
                    +------------------+-----------------------+
                                       |
                                       | Selector: app=nginx-proxy
                                       v
+------------------++------------------------------------------+
|   ConfigMap      ||          K8s Pod (nginx-proxy-xxx)       |
| (default.conf)   |+------>   Container: nginx (8080)         |
| [Proxy Rules]    ||                                          |
+------------------++---------+-------------------+------------+
                              |                   |
               / (Static HTML)|                   | /api/disk (Reverse Proxy)
                              v                   v
                    +------------------+ +-----------------------+
                    |   Static Files   | |  K8s Service          |
                    | (index.html, JS) | |  (monitor-service:8081)|
                    +------------------+ +-----------+-----------+
                                                     |
                                                     | Selector: app=monitor
                                                     v
                                         +-----------------------+
                                         |  K8s Pod (monitor-xxx)|
                                         |  Container: go (8081) |
                                         | [Disk Metrics API]    |
                                         +-----------------------+

---

## 🛠 关键实战技术 (SRE Highlights)

本项目在开发过程中解决了一系列真实的生产级挑战：

1.  **标签与选择器纠偏 (Label Selector Alignment)**: 
    解决了初始部署中存在的 Service 标签碰撞问题，确保 8080 与 8081 流量精准命中对应 Pod。
2.  **反向代理路径映射**: 
    在 Nginx ConfigMap 中实现了从根路径到 `/api/disk` 的精准转发，打通了微服务间的通信链路。
3.  **压力测试验证 (Stress Testing)**: 
    通过 `kubectl exec` 注入模拟负载（`dd` 命令生成大文件），验证了监控仪表盘的实时响应能力。

---

## 🚀 快速启动 (How to Run)

### 1. 部署到 K8s 集群
```bash
kubectl apply -f k8s/
# 转发前端端口
kubectl port-forward service/nginx-service 8080:8080 &
# 转发后端接口（用于调试）
kubectl port-forward service/monitor-service 8081:8081