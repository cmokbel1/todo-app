package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/cmokbel1/todo-app/backend/todo"
)

var (
	version = "NA"
	commit  = "NA"
	date    = "NA"
)

func main() {
	todo.Build.Version = version
	todo.Build.Commit = commit
	todo.Build.Date = date

	logger := todo.NewLogger()
	logger.SetOutput(os.Stderr)
	logger.Info(todo.BuildDetails())
	os.Exit(realMain(logger))
}

func realMain(logger todo.Logger) int {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() { <-c; cancel() }()

	logger.Info("press CTRL+C to exit")
	<-ctx.Done()
	return 0
}
