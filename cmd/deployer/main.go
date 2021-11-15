package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/crossedbot/common/golang/logger"
	"github.com/crossedbot/common/golang/server"
	"github.com/crossedbot/common/golang/service"

	"github.com/crossedbot/matryoshka/pkg/deployer/controller"
)

const (
	// Exit codes
	FatalExitCode = iota + 1
)

var (
	// Build variables
	Version = "-"
	Build   = "-"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	svc := service.New(ctx)
	err := svc.Run(run, syscall.SIGINT, syscall.SIGTERM)
	// cancel regardless of error state
	cancel()
	if err != nil {
		fatal("Error: %s", err)
	}
}

func fatal(format string, a ...interface{}) {
	logger.Error(fmt.Errorf(format, a...))
	os.Exit(FatalExitCode)
}

func newServer() server.Server {
	// TODO should have a configuration for the deployer for things like the
	// hostport
	hostport := "127.0.0.1:8080"
	srv := server.New(hostport, 30, 30)
	for _, route := range controller.Routes {
		srv.Add(
			route.Handler,
			route.Method,
			route.Path,
			route.ResponseSettings...,
		)
	}
	// Initialize the controller before returning
	controller.V1()
	return srv
}

func run(ctx context.Context) error {
	srv := newServer()
	if err := srv.Start(); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Listening on %s:%d", "127.0.0.1", 8080))
	<-ctx.Done()
	logger.Info("Received signal, shutting down...")
	return nil
}
