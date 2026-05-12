# dsu — DeepSeek Usage CLI

DeepSeek 余额查询命令行工具，自动记录每次查询快照，展示消费对比。

## 快速开始

```bash
# 一键安装
bash setup.sh

# 设置 API Key
dsu apikey sk-xxxxx

# 查询余额
dsu
```

## 命令

```
dsu                  # 余额 + 距上次查询消耗
dsu day              # 余额 + 近 24h 消耗（距上次≤24h）
dsu week             # 余额 + 近 7 天消耗（距上次≤7 天）
dsu month            # 余额 + 近 30 天消耗（距上次≤30 天）
dsu apikey <key>     # 设置 DeepSeek API Key
```

## 输出示例

```
余额:          88.09  CNY
距上次花费:      −2.15  CNY (05-10)
```

多币种：

```
余额:          88.09  CNY  |  15.50  USD
近24h花费:      −2.15  CNY  |  −0.50  USD (05-10)
```

## 手动构建

需要 Go 1.21+：

```bash
export https_proxy=http://127.0.0.1:7897  # 按需
go mod tidy
go build -o dsu .
```

## 跨平台编译

```bash
GOOS=linux GOARCH=amd64 go build -o dsu-linux .
GOOS=darwin GOARCH=arm64 go build -o dsu-mac .
GOOS=windows GOARCH=amd64 go build -o dsu.exe .
```

## 配置

API Key 存储在 `~/.dsu/data.db`（SQLite）。余额记录自动保存，每次查询写入一条。

## License

MIT
