package main

import (
	_ "go.uber.org/automaxprocs"
	"myMiniBlog/internal/miniblog"
	"os"
)

func main() {
	cmd := miniblog.NewMiniBlogCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
