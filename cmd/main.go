package main

import (
	"fmt"
	"os"

	"github.com/vincent119/timesamp/internal/cmd"
	"github.com/vincent119/timesamp/internal/i18n"
)

func main() {
	// 初始化 i18n
	if err := i18n.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize i18n: %v\n", err)
	}

	// 執行命令
	cmd.Execute()
}
