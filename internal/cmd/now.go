// Package cmd contains the command-line interface for the timestamp tool
package cmd

import (
	"fmt"
	"time"

	"timestamp/internal/converter"

	"github.com/spf13/cobra"
)

var (
	timeOffset string
)

// nowCmd 顯示當前時間
var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Show current time in various formats",
	Long: `Display current time in various formats including Unix timestamps and common date-time formats

Supported relative time offsets:
- s: seconds (e.g., +30s, -10s)
- m: minutes (e.g., +5m, -15m)  
- h: hours (e.g., +2h, -3h)
- d: days (e.g., +1d, -7d)
- w: weeks (e.g., +1w, -2w)
- M: months (e.g., +1M, -6M)
- y: years (e.g., +1y, -2y)

Examples:
  timestamp now                  # Current time
  timestamp now --offset +1d     # Tomorrow same time
  timestamp now --offset -1d     # Yesterday same time
  timestamp now --offset +1w     # Next week same time`,
	Args: cobra.NoArgs,
	RunE: showCurrentTime,
}

func init() {
	rootCmd.AddCommand(nowCmd)
	nowCmd.Flags().StringVar(&timeOffset, "offset", "", "Time offset (e.g., +1d, -1w, +2M)")
	
	// 添加 offset 參數的自動補全建議
	nowCmd.RegisterFlagCompletionFunc("offset", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			"+1s\t1 second later",
			"-1s\t1 second ago",
			"+30s\t30 seconds later",
			"-30s\t30 seconds ago",
			"+1m\t1 minute later",
			"-1m\t1 minute ago",
			"+5m\t5 minutes later",
			"-5m\t5 minutes ago",
			"+1h\t1 hour later",
			"-1h\t1 hour ago",
			"+2h\t2 hours later",
			"-2h\t2 hours ago",
			"+1d\t1 day later (tomorrow)",
			"-1d\t1 day ago (yesterday)",
			"+7d\t7 days later",
			"-7d\t7 days ago",
			"+1w\t1 week later",
			"-1w\t1 week ago",
			"+2w\t2 weeks later",
			"-2w\t2 weeks ago",
			"+1M\t1 month later",
			"-1M\t1 month ago",
			"+3M\t3 months later",
			"-3M\t3 months ago",
			"+6M\t6 months later",
			"-6M\t6 months ago",
			"+1y\t1 year later",
			"-1y\t1 year ago",
		}, cobra.ShellCompDirectiveDefault
	})
}

// showCurrentTime 顯示當前時間
func showCurrentTime(cmd *cobra.Command, args []string) error {
	// 建立轉換器
	conv, err := converter.NewConverter(timezone)
	if err != nil {
		return fmt.Errorf("failed to create converter: %v", err)
	}
	
	now := time.Now().In(conv.Location)
	
	// 處理時間偏移
	if timeOffset != "" {
		now, err = conv.AddTimeOffset(now, timeOffset)
		if err != nil {
			return fmt.Errorf("time offset calculation failed: %v", err)
		}
	}
	
	// 轉換時間
	result, err := conv.Convert(fmt.Sprintf("%d", now.Unix()), nil)
	if err != nil {
		return fmt.Errorf("conversion failed: %v", err)
	}
	
	if jsonOutput {
		outputJSON(result)
	} else {
		outputText(result)
	}
	
	return nil
}
