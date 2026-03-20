import requests
import json

# 定义你的 Go API 地址
url = "http://localhost:8081/api/disk"

try:
    # 1. 发送请求
    response = requests.get(url)
    data = response.json()
    
    # 2. 提取数据
    usage = data["percent"]
    
    # 3. 业务逻辑：自动化判断
    print(f"--- 实时监控报告 ---")
    print(f"当前磁盘占用: {usage}%")
    
    if usage > 80:
        print("⚠️  警告：磁盘空间不足！请立即清理 /tmp 目录。")
    else:
        print("✅ 状态正常：空间充足，继续运行。")

except Exception as e:
    print(f"❌ 无法连接到监控接口，请确保 Go 程序正在运行。错误: {e}")
