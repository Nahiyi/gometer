# GMeter

命令行 HTTP 压测工具，用于学习 Go 语言标准库、命令行程序以及并发知识。

## 安装

```bash
git clone https://github.com/Nahiyi/gmeter
cd gmeter
go install
```

或直接运行：

```bash
go run .
```

## 快速开始

```bash
# 1. 创建配置文件
cp req.json.example req.json
# 编辑 req.json 中的 URL 和请求参数

# 2. 运行压测
gmeter run -n 10 -t 5 -l 3 -c ./req.json -o ./report.json

# 3. 查看报告
# 用浏览器打开 viewer.html，拖入 report.json 查看
```

## 命令行参数

```bash
gmeter run <options>
```

| 参数 | 短参数 | 说明 | 默认值 |
|------|--------|------|--------|
| `--threads` | `-n` | 线程数（必填） | - |
| `--ramp-up` | `-t` | 预热时间（秒），线程渐进启动 | 0（立即启动） |
| `--loop` | `-l` | 每个线程循环次数 | 1 |
| `--config` | `-c` | 配置文件路径 | `./req.json` |
| `--output` | `-o` | JSON 报告输出路径（默认打印到终端） | stdout |
| `--request-timeout` | - | 单次请求超时（毫秒） | 5000 |
| `--max-duration` | - | 最大持续时间（秒），0 表示不限 | 0 |
| `--dry-run` | - | 验证配置格式，不实际发请求 | false |

## 配置文件

`req.json` 示例：

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
    }
  ]
}
```

- `request`：所有线程共享的请求配置

- `users`：用户隔离配置（可选）。线程 i 使用 `users[i % len(users)]`

  > 一般用于配置不同线程对应不同的用户（通过请求头的 token 区分）

- users 数组长度需要 >= 线程数

## 报告查看器

直接用浏览器打开 `viewer.html`，将 JSON 报告文件拖入即可查看可视化报表。
