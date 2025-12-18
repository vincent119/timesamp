package cmd

import (
	"os"

	"timestamp/internal/i18n"

	"github.com/spf13/cobra"
)

// completionCmd 自動補全命令
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate auto-completion scripts", // 預設英文，將在 UpdateCommandDescriptions 中更新
	Long:  "Generate auto-completion scripts", // 預設英文，將在 UpdateCommandDescriptions 中更新
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// 在 PersistentPreRun 後更新 completion 命令描述
	originalPreRun := completionCmd.PreRun
	completionCmd.PreRun = func(cmd *cobra.Command, args []string) {
		completionCmd.Short = i18n.T("cmd.completion.short")
		completionCmd.Long = i18n.T("cmd.completion.long")
		if originalPreRun != nil {
			originalPreRun(cmd, args)
		}
	}
}
