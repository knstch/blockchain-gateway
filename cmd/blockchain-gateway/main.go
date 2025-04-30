package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	defaultLog "log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
	"wallets-service/internal/wallets/blockchain"

	"github.com/knstch/subtrack-libs/endpoints"
	"github.com/knstch/subtrack-libs/log"
	"github.com/knstch/subtrack-libs/tracing"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"wallets-service/internal/endpoints/public"
	"wallets-service/internal/wallets"
	"wallets-service/internal/wallets/repo"

	"wallets-service/config"
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

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: %w", err)
	}
	dbRepo := repo.NewDBRepo(logger, db)

	blockchainRedisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
		DB:   1,
	})

	blockchainClient, err := blockchain.NewBlockchainClient(cfg, logger, blockchainRedisClient)
	if err != nil {
		return fmt.Errorf("blockchain.NewClient: %w", err)
	}

	walletsRedisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
		DB:   0,
	})

	svc := wallets.NewService(logger, dbRepo, cfg, blockchainClient, walletsRedisClient)

	publicController := public.NewController(svc, logger, cfg)
	publicEndpoints := endpoints.InitHttpEndpoints(cfg.ServiceName, publicController.Endpoints())

	srv := http.Server{
		Addr: ":" + cfg.PublicHTTPAddr,
		Handler: http.TimeoutHandler(
			publicEndpoints,
			time.Second*5,
			"service temporary unavailable",
		),
		ReadHeaderTimeout: time.Millisecond * 500,
		ReadTimeout:       time.Minute * 5,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err = srv.Shutdown(context.Background()); err != nil {
			defaultLog.Print(err)
		}
		close(idleConnsClosed)
	}()

	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-idleConnsClosed

	return nil
}
