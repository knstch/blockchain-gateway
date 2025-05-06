package main

import (
	"blockchain-gateway/config"
	"blockchain-gateway/internal/blockchain-gateway"
	"blockchain-gateway/internal/endpoints/private"
	"context"
	"fmt"
	"google.golang.org/grpc"
	defaultLog "log"
	"net"
	"os"
	"path/filepath"

	"github.com/knstch/subtrack-libs/log"
	"github.com/knstch/subtrack-libs/tracing"

	privateApi "github.com/knstch/blockchain-gateway-api/private"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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

	if err := config.InitENV(dir); err != nil {
		return fmt.Errorf("config.InitENV: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("config.GetConfig: %w", err)
	}

	shutdown := tracing.InitTracer(cfg.ServiceName, cfg.JaegerHost)
	defer shutdown(context.Background())

	logger := log.NewLogger(cfg.ServiceName, log.InfoLevel)

	svc, err := blockchain.NewService(cfg, logger)
	if err != nil {
		return fmt.Errorf("blockchain.NewService: %w", err)
	}

	privateController := private.NewController(svc, logger, cfg)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.PrivateGRPCAddr))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	privateApi.RegisterBlockchainGatewayServer(grpcServer, privateController)

	if err = grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("grpcServer.Serve: %w", err)
	}

	return nil
}
