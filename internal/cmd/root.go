// Package cmd contains the command-line interface for the timestamp tool
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"timestamp/internal/converter"

	"github.com/spf13/cobra"
)

var (
	inputFormat    string
	outputFormat   string
	timezone       string
	inputTimestamp string
	jsonOutput     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "A powerful timestamp conversion tool",
	Long: `timestamp is a feature-rich timestamp conversion tool that supports:

• Automatic detection of multiple time formats
• Unix timestamps (seconds/milliseconds/microseconds/nanoseconds)
• RFC3339, ISO8601 and other standard formats
• Custom format conversion
• Timezone support
• Relative time calculations

Examples:
  timestamp 1640995200                    # Unix timestamp conversion
  timestamp "2022-01-01 12:00:00"         # String format conversion
  timestamp -o rfc3339 1640995200         # Specify output format
  timestamp -i "2006-01-02" "2022-01-01"  # Specify input format`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			inputTimestamp = args[0]
			if err := convertTimestamp(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputFormat, "input-format", "i", "", 
		"Specify input format (unix, unix-ms, unix-us, unix-ns, rfc3339, rfc3339-nano, datetime, date, time)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output-format", "o", "datetime", 
		"Specify output format (unix, unix-ms, unix-us, unix-ns, rfc3339, rfc3339-nano, datetime, date, time)")
	rootCmd.PersistentFlags().StringVarP(&timezone, "timezone", "z", "", 
		"Specify timezone (e.g., UTC, Asia/Taipei)")
	rootCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, 
		"Output in JSON format")

	// 設定 flag 自動完成
	rootCmd.RegisterFlagCompletionFunc("timezone", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		common := []string{
			"UTC", "Local", "Asia/Taipei", "Asia/Shanghai", "Asia/Tokyo", "Asia/Seoul",
			"America/New_York", "America/Los_Angeles", "Europe/London", "Europe/Paris",
			"Australia/Sydney", "Pacific/Auckland",
		}
		return common, cobra.ShellCompDirectiveDefault
	})

	rootCmd.RegisterFlagCompletionFunc("input-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		formats := []string{
			"unix", "unix-ms", "unix-us", "unix-ns",
			"rfc3339", "rfc3339-nano", "datetime", "date", "time",
		}
		return formats, cobra.ShellCompDirectiveDefault
	})

	rootCmd.RegisterFlagCompletionFunc("output-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		formats := []string{
			"unix", "unix-ms", "unix-us", "unix-ns",
			"rfc3339", "rfc3339-nano", "datetime", "date", "time",
		}
		return formats, cobra.ShellCompDirectiveDefault
	})
}

func convertTimestamp(cmd *cobra.Command, args []string) error {
	conv, err := converter.NewConverter(timezone)
	if err != nil {
		return fmt.Errorf("failed to create converter: %v", err)
	}

	// 解析輸入格式
	var inputFmt *converter.TimestampFormat
	if inputFormat != "" {
		format, parseErr := parseInputFormat(inputFormat)
		if parseErr != nil {
			return parseErr
		}
		inputFmt = &format
	}

	// 轉換時間戳
	result, err := conv.Convert(inputTimestamp, inputFmt)
	if err != nil {
		return fmt.Errorf("conversion failed: %v", err)
	}

	// 輸出結果
	if jsonOutput {
		outputJSON(result)
	} else {
		outputText(result)
	}
	return nil
}

func parseInputFormat(format string) (converter.TimestampFormat, error) {
	switch format {
	case "unix":
		return converter.UnixSeconds, nil
	case "unix-ms":
		return converter.UnixMilliseconds, nil
	case "unix-us":
		return converter.UnixMicroseconds, nil
	case "unix-ns":
		return converter.UnixNanoseconds, nil
	case "rfc3339":
		return converter.RFC3339, nil
	case "rfc3339-nano":
		return converter.RFC3339Nano, nil
	case "datetime":
		return converter.DateTime, nil
	case "date":
		return converter.DateOnly, nil
	case "time":
		return converter.TimeOnly, nil
	default:
		return 0, fmt.Errorf("unsupported format: %s", format)
	}
}

func outputText(result *converter.ConvertResult) {
	fmt.Printf("Original Input: %s\n", result.Original)
	fmt.Printf("Detected Format: %s\n", result.DetectedFormat)
	
	// 根據輸出格式選擇顯示
	switch outputFormat {
	case "unix":
		fmt.Printf("Converted: %d\n", result.UnixSeconds)
	case "unix-ms":
		fmt.Printf("Converted: %d\n", result.UnixMillis)
	case "unix-us":
		fmt.Printf("Converted: %d\n", result.UnixMicros)
	case "unix-ns":
		fmt.Printf("Converted: %d\n", result.UnixNanos)
	case "rfc3339":
		fmt.Printf("Converted: %s\n", result.RFC3339)
	case "rfc3339-nano":
		fmt.Printf("Converted: %s\n", result.RFC3339Nano)
	case "date":
		fmt.Printf("Converted: %s\n", result.DateOnly)
	case "time":
		fmt.Printf("Converted: %s\n", result.TimeOnly)
	default: // datetime
		fmt.Printf("Converted: %s\n", result.DateTime)
	}
	
	fmt.Printf("Unix Timestamp: %d\n", result.UnixSeconds)
	fmt.Printf("Weekday: %s\n", result.Weekday)
	fmt.Printf("Timezone: %s\n", result.Timezone)
}

func outputJSON(result *converter.ConvertResult) {
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
}
