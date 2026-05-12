package main

import (
	"fmt"
	"os"

	"dsu/cmd"
	"dsu/store"
)

func main() {
	if err := store.Init(); err != nil {
		fmt.Fprintln(os.Stderr, "初始化失败:", err)
		os.Exit(1)
	}
	defer store.Close()

	cmd.Execute()
}
