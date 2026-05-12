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
	Args:  cobra.ArbitraryArgs,
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

		if err := store.InsertRecords(current); err != nil {
			fmt.Fprintln(os.Stderr, "数据库错误:", err)
			os.Exit(1)
		}

		switch window {
		case "":
			lastRecords, lastTime, _ := store.LastRecords()
			display.Print(current, recordsToMap(lastRecords), lastTime, "距上次花费")
		case "day":
			showWindow(current, 24*time.Hour, "近24h花费")
		case "week":
			showWindow(current, 7*24*time.Hour, "近7天花费")
		case "month":
			showWindow(current, 30*24*time.Hour, "近30天花费")
		default:
			lastRecords, lastTime, _ := store.LastRecords()
			display.Print(current, recordsToMap(lastRecords), lastTime, "距上次花费")
		}
	},
}

func recordsToMap(records []store.Record) map[string]float64 {
	m := make(map[string]float64)
	for _, r := range records {
		m[r.Currency] = r.TotalBalance
	}
	return m
}

func showWindow(current map[string]float64, duration time.Duration, label string) {
	records, err := store.RecordsSince(duration)
	if err != nil {
		fmt.Fprintln(os.Stderr, "数据库错误:", err)
		os.Exit(1)
	}

	if len(records) == 0 {
		display.Print(current, nil, "", "")
		return
	}

	earliest := make(map[string]float64)
	var earliestTime string
	for _, r := range records {
		if _, ok := earliest[r.Currency]; !ok {
			earliest[r.Currency] = r.TotalBalance
		}
		if earliestTime == "" || r.CreatedAt < earliestTime {
			earliestTime = r.CreatedAt
		}
	}

	display.Print(current, earliest, earliestTime, label)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
