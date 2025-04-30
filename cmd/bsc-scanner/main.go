package main

import (
	"blockchain-gateway/config"
	"blockchain-gateway/internal/infra/ethnode"
	"context"
	"fmt"
	"github.com/knstch/subtrack-libs/log"
	"github.com/knstch/subtrack-libs/tracing"
	defaultLog "log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		defaultLog.Println(err)
		recover()
	}
}

func run() error {
	args := os.Args

	dir, err := filepath.Abs(filepath.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}

	if err = config.InitENV(dir); err != nil {
		return fmt.Errorf("config.InitENV: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("config.GetConfig: %w", err)
	}

	shutdown := tracing.InitTracer(cfg.ServiceName, cfg.JaegerHost)
	defer shutdown(context.Background())

	logger := log.NewLogger(cfg.ServiceName, log.InfoLevel)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	node, err := ethnode.MakeBscMempoolScanner(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("ethnode.MakeBscMempoolScanner: %w", err)
	}

	if err = node.Scan(ctx); err != nil {
		return fmt.Errorf("node.Scan: %w", err)
	}

	return nil
}
