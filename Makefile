# Makefile for timestamp CLI tool

.PHONY: build clean test install help run demo build-all build-windows build-macos build-linux

# 預設目標
all: build

# 編譯
build:
	@echo "編譯 timestamp..."
	@go build -o timestamp cmd/main.go
	@echo "編譯完成！執行檔: ./timestamp"

# 跨平台編譯
build-all: build-windows build-macos build-linux
	@echo "所有平台編譯完成！"

# Windows 編譯
build-windows:
	@echo "編譯 Windows 版本..."
	@mkdir -p dist
	@GOOS=windows GOARCH=amd64 go build -o dist/timestamp-windows-amd64.exe cmd/main.go
	@GOOS=windows GOARCH=386 go build -o dist/timestamp-windows-386.exe cmd/main.go
	@GOOS=windows GOARCH=arm64 go build -o dist/timestamp-windows-arm64.exe cmd/main.go
	@echo "Windows 版本編譯完成！"

# macOS 編譯
build-macos:
	@echo "編譯 macOS 版本..."
	@mkdir -p dist
	@GOOS=darwin GOARCH=amd64 go build -o dist/timestamp-darwin-amd64 cmd/main.go
	@GOOS=darwin GOARCH=arm64 go build -o dist/timestamp-darwin-arm64 cmd/main.go
	@echo "macOS 版本編譯完成！"

# Linux 編譯
build-linux:
	@echo "編譯 Linux 版本..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -o dist/timestamp-linux-amd64 cmd/main.go
	@GOOS=linux GOARCH=386 go build -o dist/timestamp-linux-386 cmd/main.go
	@GOOS=linux GOARCH=arm64 go build -o dist/timestamp-linux-arm64 cmd/main.go
	@GOOS=linux GOARCH=arm go build -o dist/timestamp-linux-arm cmd/main.go
	@echo "Linux 版本編譯完成！"

# 編譯並安裝到 GOPATH/bin
install:
	@echo "安裝 timestamp..."
	@go install ./cmd
	@echo "安裝完成！可直接使用 timestamp 命令"

# 清理編譯產物
clean:
	@echo "清理編譯產物..."
	@rm -f timestamp
	@rm -rf dist/
	@rm -f _timestamp timestamp-completion.bash timestamp.fish timestamp.ps1
	@echo "清理完成"

# 執行測試
test:
	@echo "執行測試..."
	@go test ./...
	@echo "測試完成"

# 整理依賴
tidy:
	@echo "整理依賴..."
	@go mod tidy
	@echo "依賴整理完成"

# 格式化程式碼
fmt:
	@echo "格式化程式碼..."
	@go fmt ./...
	@echo "格式化完成"

# 檢查程式碼
lint:
	@echo "檢查程式碼..."
	@go vet ./...
	@echo "檢查完成"

# 生成自動補全腳本
completion: build
	@echo "生成自動補全腳本..."
	@mkdir -p completions
	@./timestamp completion bash > completions/timestamp-completion.bash
	@./timestamp completion zsh > completions/_timestamp
	@./timestamp completion fish > completions/timestamp.fish
	@./timestamp completion powershell > completions/timestamp.ps1
	@echo "自動補全腳本已生成到 completions/ 目錄"

# 安裝 zsh 自動補全 (macOS)
install-zsh-completion: build
	@echo "安裝 zsh 自動補全..."
	@./timestamp completion zsh > _timestamp
	@if [ -d "/usr/local/share/zsh/site-functions" ]; then \
		sudo mv _timestamp /usr/local/share/zsh/site-functions/; \
		echo "已安裝到系統目錄: /usr/local/share/zsh/site-functions/"; \
	elif [ -d "$$HOME/.oh-my-zsh/completions" ]; then \
		mv _timestamp $$HOME/.oh-my-zsh/completions/; \
		echo "已安裝到 oh-my-zsh: $$HOME/.oh-my-zsh/completions/"; \
	else \
		mkdir -p $$HOME/.local/share/zsh/site-functions; \
		mv _timestamp $$HOME/.local/share/zsh/site-functions/; \
		echo "已安裝到用戶目錄: $$HOME/.local/share/zsh/site-functions/"; \
	fi
	@echo "請重新載入 zsh 或執行: source ~/.zshrc"

# 執行示範
demo: build
	@echo "=== 時間戳轉換工具示範 ==="
	@echo
	@echo "1. 轉換 Unix 時間戳:"
	@./timestamp 1642781234
	@echo
	@echo "2. 轉換日期時間:"
	@./timestamp "2022-01-21 12:00:34"
	@echo
	@echo "3. 顯示當前時間:"
	@./timestamp now
	@echo
	@echo "4. 顯示明天時間:"
	@./timestamp now --offset +1d
	@echo
	@echo "5. 顯示上週時間:"
	@./timestamp now --offset -1w
	@echo
	@echo "6. JSON 格式輸出:"
	@./timestamp 1642781234 --json

# 顯示幫助
help:
	@echo "可用的 make 目標:"
	@echo "  build              - 編譯程式 (當前平台)"
	@echo "  build-all          - 編譯所有平台版本"
	@echo "  build-windows      - 編譯 Windows 版本"
	@echo "  build-macos        - 編譯 macOS 版本"
	@echo "  build-linux        - 編譯 Linux 版本"
	@echo "  install            - 編譯並安裝到 GOPATH/bin"
	@echo "  clean              - 清理編譯產物和補全腳本"
	@echo "  test               - 執行測試"
	@echo "  tidy               - 整理 Go 模組依賴"
	@echo "  fmt                - 格式化程式碼"
	@echo "  lint               - 檢查程式碼"
	@echo "  completion         - 生成所有 shell 的自動補全腳本"
	@echo "  install-zsh-completion - 安裝 zsh 自動補全"
	@echo "  demo               - 執行示範"
	@echo "  help               - 顯示此幫助訊息"
