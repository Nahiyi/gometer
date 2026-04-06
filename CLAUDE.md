# GoMeter

命令行 HTTP 压测/接口探测工具，用于学习 Go 语言标准库（HTTP、flag、并发控制等）。

## 工具信息

- **项目名/工具名**: `gometer`
- **定位**: 压测工具 / HTTP 接口探测客户端
- **Go 版本**: 1.25.1

## 命令行接口

```bash
gometer run [参数]

核心参数：
  -n, --threads        线程数（必填）
  -t, --ramp-up        Ramp-Up 秒数（默认 0，立即启动）
  -l, --loop           每个线程循环次数（默认 1）
  -c, --config         请求配置文件（默认 ./req.json）
  -o, --output         JSON 报告输出路径（默认 stdout）

可选参数：
  --request-timeout    单次请求超时 ms（默认 5000）
  --max-duration       最大持续时间 秒（默认 0，不限）
  --dry-run            验证配置格式，不实际发请求
```

## 配置文件结构

配置文件为 JSON 格式，示例见 `req.json.example`。

```json
{
  "request": {
    "url": "http://example.com/api",
    "method": "POST",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{}"
  },
  "users": [
    {
      "headers": {
        "token": "user-token-1",
        "x-user-id": "1001"
      }
    },
    {
      "headers": {
        "token": "user-token-2",
        "x-user-id": "1002"
      }
    }
  ]
}
```

**说明：**
- `request`：所有线程共享的请求配置（URL、Method、Headers、Body）
- `users`：用户隔离配置，数组长度 >= 线程数，线程 i 使用 users[i]
- `users[].headers`：请求头隔离字段，会与 `request.headers` 合并
- users 数组元素预留扩展，未来可加 `queryParams`、`body` 等字段

## 核心行为

### 线程启动（Ramp-Up）
- `ramp-up=0`：所有线程瞬间启动
- `ramp-up=10, 线程数=100`：每 0.1 秒启动一个线程（10s / 100）
- Ramp-Up 只管启动，不管请求之间的间隔

### 循环执行
- 每个线程按 `loop` 次循环执行请求
- 同一个线程的多次循环，使用同一个 user 配置
- 循环之间无额外延迟

### 错误处理
- 请求失败（如超时、500、400 等）记录错误信息，继续执行
- 成功和失败都是合法结果，都记录到报告中

## JSON 报告结构

```json
{
  "summary": {
    "totalThreads": 100,
    "totalLoops": 10,
    "totalRequests": 1000,
    "successCount": 980,
    "failCount": 20,
    "successRate": 0.98,
    "durationMs": 15230,
    "avgResponseTimeMs": 125.5,
    "minResponseTimeMs": 45,
    "maxResponseTimeMs": 3500,
    "p50ResponseTimeMs": 110,
    "p90ResponseTimeMs": 200,
    "p99ResponseTimeMs": 500
  },
  "threads": [
    {
      "threadId": 1,
      "loopResults": [
        {
          "loopIndex": 1,
          "requests": [
            {
              "requestIndex": 1,
              "url": "http://example.com/api",
              "method": "POST",
              "requestHeaders": {...},
              "requestBody": "{}",
              "responseStatus": 200,
              "responseTimeMs": 120,
              "responseHeaders": {...},
              "success": true,
              "error": null
            },
            {
              "requestIndex": 2,
              "responseStatus": 500,
              "responseTimeMs": 50,
              "success": false,
              "error": "server error"
            }
          ]
        }
      ]
    }
  ]
}
```

## 模块划分

```
gometer/
├── cmd/
│   ├── root.go           # CLI 入口，flag 定义
│   └── run.go            # run 子命令
├── internal/
│   ├── config/           # 配置文件解析（JSON -> struct）
│   ├── loader/            # 用户数据加载（users 数组 -> 线程映射）
│   ├── httpclient/       # HTTP 客户端封装（超时控制、请求发送）
│   ├── runner/            # 压测运行器（线程调度、ramp-up、循环执行）
│   ├── collector/         # 结果收集（并发安全的结果聚合）
│   └── reporter/         # JSON 报告生成
├── req.json.example      # 示例配置文件
└── CLAUDE.md
```

## 设计原则

1. **学习优先** - 代码规范清晰、模块分明确、不过度设计
2. **最佳实践** - 遵循 Go 语言惯用模式
3. **YAGNI** - 不做过度抽象，能简单实现就不复杂化
4. **错误即结果** - 请求失败也是一种合法结果，都需要记录
