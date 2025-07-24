# Timestamp 轉換工具

一個使用 Go 語言和 Cobra 框架開發的時間戳轉換 CLI 工具，支援多種時間格式間的相互轉換。

## 功能特色

- 🚀 自動偵測輸入的時間格式
- 🌍 支援時區轉換
- 📊 多種輸出格式 (預設、JSON、指定格式)
- ⚡ 快速且輕量的 CLI 工具
- 🎯 支援豐富的時間格式

## 支援的格式

- Unix 時間戳 (秒、毫秒、微秒、納秒)
- RFC3339 格式
- RFC3339Nano 格式
- 日期時間格式 (YYYY-MM-DD HH:MM:SS)
- 日期格式 (YYYY-MM-DD)
- 時間格式 (HH:MM:SS)

## 安裝

### 從原始碼編譯

```bash
git clone https://github.com/vincent119/timesamp.git
cd timestamp
go build -o timestamp cmd/main.go
```

### 跨平台編譯

#### Windows

```bash
# Windows 64位
GOOS=windows GOARCH=amd64 go build -o timestamp.exe cmd/main.go

# Windows 32位
GOOS=windows GOARCH=386 go build -o timestamp.exe cmd/main.go

# Windows ARM64
GOOS=windows GOARCH=arm64 go build -o timestamp.exe cmd/main.go
```

#### macOS

```bash
# macOS Intel (x86_64)
GOOS=darwin GOARCH=amd64 go build -o timestamp cmd/main.go

# macOS Apple Silicon (ARM64)
GOOS=darwin GOARCH=arm64 go build -o timestamp cmd/main.go

# 通用二進制檔案 (同時支援 Intel 和 Apple Silicon)
# 需要分別編譯後合併
GOOS=darwin GOARCH=amd64 go build -o timestamp-amd64 cmd/main.go
GOOS=darwin GOARCH=arm64 go build -o timestamp-arm64 cmd/main.go
lipo -create -output timestamp timestamp-amd64 timestamp-arm64
```

#### Linux

```bash
# Linux 64位
GOOS=linux GOARCH=amd64 go build -o timestamp cmd/main.go

# Linux 32位
GOOS=linux GOARCH=386 go build -o timestamp cmd/main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o timestamp cmd/main.go

# Linux ARM (Raspberry Pi 等)
GOOS=linux GOARCH=arm go build -o timestamp cmd/main.go
```

#### 一鍵編譯所有平台

您也可以使用 Makefile 來編譯所有平台：

```bash
# 編譯所有平台
make build-all

# 或者單獨編譯特定平台
make build-windows
make build-macos
make build-linux
```

## 自動補全設定

工具支援 Bash、Zsh、Fish 和 PowerShell 的自動補全功能。

### Zsh 自動補全

#### 方法一：臨時啟用 (當前 session)

```bash
# 直接載入到當前 session
source <(./timestamp completion zsh)
```

#### 方法二：永久安裝 (推薦)

```bash
# 生成補全腳本
./timestamp completion zsh > _timestamp

# 安裝到系統目錄 (需要 admin 權限)
sudo mv _timestamp /usr/local/share/zsh/site-functions/

# 或者安裝到用戶目錄
mkdir -p ~/.local/share/zsh/site-functions
mv _timestamp ~/.local/share/zsh/site-functions/

# 如果使用 oh-my-zsh
mkdir -p ~/.oh-my-zsh/completions
./timestamp completion zsh > ~/.oh-my-zsh/completions/_timestamp
```

#### 方法三：添加到 ~/.zshrc

```bash
# 添加到 zsh 配置檔案
echo 'source <(timestamp completion zsh)' >> ~/.zshrc

# 重新載入配置
source ~/.zshrc
```

### Bash 自動補全

#### macOS (使用 Homebrew)

```bash
# 安裝 bash-completion (如果尚未安裝)
brew install bash-completion

# 生成並安裝補全腳本
timestamp completion bash > /usr/local/etc/bash_completion.d/timestamp

# 重新載入 bash
source ~/.bash_profile
```

#### Linux

```bash
# 生成補全腳本
timestamp completion bash > timestamp-completion.bash

# 安裝到系統目錄
sudo mv timestamp-completion.bash /etc/bash_completion.d/

# 或者添加到 ~/.bashrc
echo 'source <(timestamp completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Fish 自動補全

```bash
# 生成並安裝補全腳本
timestamp completion fish > ~/.config/fish/completions/timestamp.fish

# 重新載入 fish
fish -c "source ~/.config/fish/completions/timestamp.fish"
```

### PowerShell 自動補全 (Windows)

```powershell
# 生成補全腳本
./timestamp.exe completion powershell > timestamp.ps1

# 添加到 PowerShell Profile
Add-Content $PROFILE '. .\timestamp.ps1'

# 重新載入 Profile
. $PROFILE
```

### 驗證自動補全

安裝完成後，您可以測試自動補全功能：

```bash
# 按 Tab 鍵查看可用命令
timestamp <Tab>

# 按 Tab 鍵查看 timezone 選項
timestamp --timezone <Tab>

# 按 Tab 鍵查看 now 命令的 offset 選項
timestamp now --offset <Tab>
```

## 使用方法

### 基本用法

```bash
# 轉換 Unix 時間戳
./timestamp 1642781234

# 轉換日期時間
./timestamp "2022-01-21 12:00:34"

# 轉換 RFC3339 格式
./timestamp "2022-01-21T12:00:34Z"

# 轉換日期
./timestamp "2022-01-21"

# 轉換時間
./timestamp "12:00:34"
```

### 選項參數

```bash
# 指定時區
./timestamp 1642781234 --timezone "UTC"
./timestamp 1642781234 --timezone "Asia/Taipei"

# 指定輸入格式
./timestamp 1642781234 --input-format unix-s

# 指定輸出格式
./timestamp 1642781234 --output-format unix-ms

# JSON 格式輸出
./timestamp 1642781234 --json
```

### 子命令

```bash
# 顯示當前時間
./timestamp now

# 顯示當前時間 (JSON 格式)
./timestamp now --json

# 顯示當前時間 (指定時區)
./timestamp now --timezone "UTC"

# 相對時間偏移
./timestamp now --offset +1d      # 明天同一時間
./timestamp now --offset -1d      # 昨天同一時間
./timestamp now --offset +1w      # 下週同一時間
./timestamp now --offset -1w      # 上週同一時間
./timestamp now --offset +1M      # 下個月同一時間
./timestamp now --offset -1M      # 上個月同一時間

# 縮寫形式
./timestamp now -o +1d            # 明天
./timestamp now -o -1w            # 上週
```

## 範例

### 基本轉換

```bash
$ ./timestamp 1642781234
原始輸入: 1642781234
偵測格式: Unix 秒級時間戳
時區: Local (CST, UTC+08:00)
星期: 星期六

轉換結果:
  Unix 秒級時間戳:   1642781234
  Unix 毫秒級時間戳: 1642781234000
  Unix 微秒級時間戳: 1642781234000000
  Unix 納秒級時間戳: 1642781234000000000
  RFC3339 格式:      2022-01-22T00:07:14+08:00
  RFC3339Nano 格式:  2022-01-22T00:07:14+08:00
  日期時間格式:      2022-01-22 00:07:14
  日期格式:          2022-01-22
  時間格式:          00:07:14
```

### JSON 輸出

```bash
$ ./timestamp 1642781234 --json
{
  "original": "1642781234",
  "detected_format": "Unix 秒級時間戳",
  "unix_seconds": 1642781234,
  "unix_milliseconds": 1642781234000,
  "unix_microseconds": 1642781234000000,
  "unix_nanoseconds": 1642781234000000000,
  "rfc3339": "2022-01-22T00:07:14+08:00",
  "rfc3339_nano": "2022-01-22T00:07:14+08:00",
  "datetime": "2022-01-22 00:07:14",
  "date_only": "2022-01-22",
  "time_only": "00:07:14",
  "weekday": "星期六",
  "timezone": "Local (CST, UTC+08:00)"
}
```

### 指定輸出格式

```bash
$ ./timestamp 1642781234 --output-format unix-ms
1642781234000

$ ./timestamp "2022-01-21 12:00:34" --output-format rfc3339
2022-01-21T12:00:34+08:00
```

## 專案結構

```
timestamp/
├── cmd/
│   └── main.go              # 主程式入口
├── internal/
│   ├── cmd/
│   │   ├── root.go          # Cobra 根命令
│   │   └── now.go           # now 子命令
│   └── converter/
│       └── converter.go     # 時間轉換核心邏輯
├── go.mod
├── go.sum
└── README.md
```

本專案採用 Go 官方建議的目錄結構：

- `cmd/`: 包含應用程式的主要入口點
- `internal/`: 包含私有的應用程式和函式庫程式碼
- `internal/cmd/`: CLI 命令實作
- `internal/converter/`: 時間轉換邏輯

## 支援的相對時間偏移

工具支援以下時間單位的偏移：

| 單位 | 說明 | 範例           |
| ---- | ---- | -------------- |
| `s`  | 秒   | `+30s`, `-10s` |
| `m`  | 分鐘 | `+5m`, `-15m`  |
| `h`  | 小時 | `+2h`, `-3h`   |
| `d`  | 天   | `+1d`, `-7d`   |
| `w`  | 週   | `+1w`, `-2w`   |
| `M`  | 月   | `+1M`, `-6M`   |
| `y`  | 年   | `+1y`, `-2y`   |

## 支援的輸入/輸出格式標識

| 格式              | 標識           | 範例                             |
| ----------------- | -------------- | -------------------------------- |
| Unix 秒級時間戳   | `unix-s`       | `1642781234`                     |
| Unix 毫秒級時間戳 | `unix-ms`      | `1642781234000`                  |
| Unix 微秒級時間戳 | `unix-us`      | `1642781234000000`               |
| Unix 納秒級時間戳 | `unix-ns`      | `1642781234000000000`            |
| RFC3339           | `rfc3339`      | `2022-01-21T12:00:34Z`           |
| RFC3339Nano       | `rfc3339-nano` | `2022-01-21T12:00:34.123456789Z` |
| 日期時間          | `datetime`     | `2022-01-21 12:00:34`            |
| 日期              | `date`         | `2022-01-21`                     |
| 時間              | `time`         | `12:00:34`                       |

## 時區支援

工具支援所有標準時區，包括但不限於：

- `Local` (系統本地時區)
- `UTC`
- `Asia/Taipei`
- `America/New_York`
- `Europe/London`
- `Asia/Tokyo`

## 依賴

- [Cobra](https://github.com/spf13/cobra) - 強大的 CLI 框架

## 開發

### 編譯

```bash
go build -o timestamp cmd/main.go
```

### 測試

```bash
go test ./...
```

### 安裝依賴

```bash
go mod tidy
```

## 授權

MIT License
