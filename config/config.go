package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"path/filepath"
)

type Config struct {
	JwtSecret string `envconfig:"JWT_SECRET" required:"true"`

	JaegerHost  string `envconfig:"JAEGER_HOST" required:"true"`
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`

	PublicHTTPAddr  string `envconfig:"PUBLIC_HTTP_ADDR" required:"true"`
	PrivateGRPCAddr string `envconfig:"PRIVATE_GRPC_ADDR" required:"true"`

	Blockchains BlockchainConfig
}

type BlockchainConfig struct {
	BscAddr string `envconfig:"BSC_ADDR" required:"true"`
}

func GetConfig() (*Config, error) {
	config := &Config{}

	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func InitENV(dir string) error {
	if err := godotenv.Load(filepath.Join(dir, ".env.local")); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("godotenv.Load: %w", err)
		}
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		return fmt.Errorf("godotenv.Load: %w", err)
	}
	return nil
}
