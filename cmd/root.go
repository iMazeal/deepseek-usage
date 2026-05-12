package cmd

import (
	"fmt"
	"os"
	"time"

	"dsu/api"
	"dsu/display"
	"dsu/store"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dsu [day|week|month]",
	Short: "查询 DeepSeek 余额和消费",
	Run: func(cmd *cobra.Command, args []string) {
		window := ""
		if len(args) > 0 {
			window = args[0]
		}

		apikey, err := store.GetConfig("apikey")
		if err != nil {
			fmt.Fprintln(os.Stderr, "数据库错误:", err)
			os.Exit(1)
		}
		if apikey == "" {
			fmt.Fprintln(os.Stderr, "请先设置: dsu apikey <key>")
			os.Exit(1)
		}

		current, err := api.FetchBalance(apikey)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		lastRecords, lastTime, err := store.LastRecords()
		if err != nil {
			fmt.Fprintln(os.Stderr, "数据库错误:", err)
			os.Exit(1)
		}

		if err := store.InsertRecords(current); err != nil {
			fmt.Fprintln(os.Stderr, "数据库错误:", err)
			os.Exit(1)
		}

		lastMap := recordsToMap(lastRecords)

		label, showConsumption := resolveWindow(window, lastTime)
		if !showConsumption {
			lastMap = nil
			lastTime = ""
		}

		display.Print(current, lastMap, lastTime, label)
	},
}

func recordsToMap(records []store.Record) map[string]float64 {
	m := make(map[string]float64)
	for _, r := range records {
		m[r.Currency] = r.TotalBalance
	}
	return m
}

func resolveWindow(window, lastTime string) (label string, show bool) {
	if window == "" {
		return "距上次花费", true
	}

	if lastTime == "" {
		switch window {
		case "day":
			return "近24h花费", false
		case "week":
			return "近7天花费", false
		case "month":
			return "近30天花费", false
		}
	}

	t, err := time.Parse(time.RFC3339, lastTime)
	if err != nil {
		return "距上次花费", false
	}
	elapsed := time.Since(t)

	switch window {
	case "day":
		if elapsed <= 24*time.Hour {
			return "近24h花费", true
		}
		return "近24h花费", false
	case "week":
		if elapsed <= 7*24*time.Hour {
			return "近7天花费", true
		}
		return "近7天花费", false
	case "month":
		if elapsed <= 30*24*time.Hour {
			return "近30天花费", true
		}
		return "近30天花费", false
	default:
		return "距上次花费", true
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
