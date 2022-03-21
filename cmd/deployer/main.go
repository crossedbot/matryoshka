package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/crossedbot/common/golang/config"
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

type Config struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	ReadTimeout  int    `toml:"read_timeout"`
	WriteTimeout int    `toml:"write_timeout"`
}

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

func newServer(c Config) server.Server {
	hostport := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	srv := server.New(hostport, c.ReadTimeout, c.WriteTimeout)
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
	f := ParseFlags()
	if f.Version {
		fmt.Printf(
			"%s version %s, build %s\n",
			filepath.Base(os.Args[0]),
			Version, Build,
		)
		return nil
	}
	config.Path(f.ConfigFile)
	var c Config
	if err := config.Load(&c); err != nil {
		return err
	}
	srv := newServer(c)
	if err := srv.Start(); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Listening on %s:%d", c.Host, c.Port))
	<-ctx.Done()
	logger.Info("Received signal, shutting down...")
	return nil
}
