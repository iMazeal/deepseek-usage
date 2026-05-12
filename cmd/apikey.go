package cmd

import (
	"fmt"
	"os"

	"dsu/store"

	"github.com/spf13/cobra"
)

var apikeyCmd = &cobra.Command{
	Use:   "apikey <key>",
	Short: "设置 DeepSeek API Key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := store.SetConfig("apikey", args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "设置失败:", err)
			os.Exit(1)
		}
		fmt.Println("API Key 已设置")
	},
}

func init() {
	rootCmd.AddCommand(apikeyCmd)
}
