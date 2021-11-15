package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/crossedbot/common/golang/logger"
	"github.com/crossedbot/common/golang/service"

	"github.com/crossedbot/matryoshka/pkg/runner"
)

const (
	DefaultOnce   = true
	FatalExitCode = iota + 1
)

var (
	// Build variables
	Version = "-"
	Build   = "-"
)

func main() {
	ctx := context.Background()
	svc := service.New(ctx)
	if err := svc.Run(run, syscall.SIGINT, syscall.SIGTERM); err != nil {
		fatal("Error: %s", err)
	}
}

func fatal(format string, a ...interface{}) {
	logger.Error(fmt.Errorf(format, a...))
	os.Exit(FatalExitCode)
}

func run(ctx context.Context) error {
	return runner.New(ctx).Run(DefaultOnce)
}
