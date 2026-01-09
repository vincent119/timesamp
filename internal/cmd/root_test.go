package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func executeCommand(args ...string) (string, error) {
	// 重置 flag 以避免測試間干擾
	rootCmd.SetArgs(args)

	// 捕獲輸出
	b := bytes.NewBufferString("")

	// 保存原本的 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 執行後恢復 stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	// 執行命令，同時捕獲 rootCmd.Execute 可能產生的錯誤 (雖然 Execute 本身會 exit，但我們可以測試 command 直接 Run)
	// 這裡我們直接呼叫 rootCmd.ExecuteC 來獲得更細緻的控制，或者為了測試可測試性，我們只測試 Run

	// 因為 rootCmd.Execute() 會在錯誤時 os.Exit(1)，我們需要小心。
	// 更好的方式是直接測試 rootCmd.Run 或類似邏輯，或者使用 SetOutput 與 SetArgs 並呼叫 ExecuteC

	rootCmd.SetOut(b)
	rootCmd.SetErr(b)

	// 因為 Execute 內部有 os.Exit(1) 在錯誤時，我們可能無法完全測試錯誤情況而不 crash process。
	// 但我們可以測試成功的情況。
	// 對於 root command，Execute() 是入口。

	// 注意：root.go 中的 Execute() 函式會呼叫 os.Exit。我們應該直接呼叫 rootCmd.Execute() 來測試，並捕獲 error
	// 但是 rootCmd 的 Run 函式也有可能 os.Exit(1) 如果 convertTimestamp 失敗。
	// 在 root.go:48: os.Exit(1)。這對測試不友善。
	// 我們可以測試那些不會導致 os.Exit 的路徑。

	err := rootCmd.Execute()

	// 讀取 pipe 內容
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// 合併 buffer
	output := b.String() + buf.String()

	return output, err
}

// 由於 root.go:43 的 Run 函式中，錯誤會導致 os.Exit(1)，這使得單元測試很難撰寫。
// 為了測試，我們可以直接測試 convertTimestamp 邏輯，或測試那些成功的案例。
// 或者，我們可以修改 root.go 讓它更易測試（例如傳回 error 而不是 exit），但根據規則我應儘量不修改原有邏輯除非必要。
// 不過為了增加測試，微調 root.go 結構是合理的。
// 但現在我先測試成功的路徑。

func TestRootCmd_UnixTimestamp(t *testing.T) {
	// 測試 Unix 時間戳轉換
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "Convert Unix timestamp",
			args: []string{"1642781234"},
			contains: []string{
				"Original Input: 1642781234",
				"Detected Format: Unix 秒級時間戳", // 假設這是在中文環境下，或預設環境
				"Converted: 2022-01-2", // 部分匹配日期
			},
		},
		{
			name: "Convert with output format",
			args: []string{"1642781234", "--output-format", "unix-ms"},
			contains: []string{
				"Converted: 1642781234000",
			},
		},
		{
			name: "Convert with JSON output",
			args: []string{"1642781234", "--json"},
			contains: []string{
				"\"original\": \"1642781234\"",
				"\"unix_seconds\": 1642781234",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 捕獲輸出
			output, err := executeCommandC(tt.args...)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			for _, c := range tt.contains {
				if !strings.Contains(output, c) {
					t.Errorf("Output should contain %q, got:\n%s", c, output)
				}
			}
		})
	}
}

func TestNowCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "Now command basic",
			args: []string{"now"},
			contains: []string{
				"Original Input:",
				"Detected Format:",
				"Unix Timestamp:",
			},
		},
		{
			name: "Now command with json",
			args: []string{"now", "--json"},
			contains: []string{
				"\"unix_seconds\":",
				"\"datetime\":",
			},
		},
		{
			name: "Now command with offset",
			args: []string{"now", "--offset", "+1h"},
			contains: []string{
				"Converted:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommandC(tt.args...)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			for _, c := range tt.contains {
				if !strings.Contains(output, c) {
					t.Errorf("Output should contain %q, got:\n%s", c, output)
				}
			}
		})
	}
}


// executeCommandC 是一個輔助函數，用來正確地執行 Cobra 命令並捕獲輸出
func executeCommandC(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	// 重置 flags 以避免測試污染 (特別是 persistent flags)
	// rootCmd.ResetFlags() // 這樣會清除所有 flags 定義，這不是我們想要的
	// 相反，我們需要手動重置 var 變數，這些變數綁定到 flags
	resetFlags()

	// 捕獲 stdout (因為 outputText 使用 fmt.Printf)
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	return string(out) + buf.String(), err
}

func resetFlags() {
	inputFormat = ""
	outputFormat = "datetime" // 預設值
	timezone = ""
	inputTimestamp = ""
	jsonOutput = false
	langFlag = ""
	timeOffset = ""
}
