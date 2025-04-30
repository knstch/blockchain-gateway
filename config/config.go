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

	WalletSecret string `envconfig:"WALLET_SECRET" required:"true"`

	KafkaAddr string `envconfig:"KAFKA_ADDR" required:"true"`

	DBConfig DBConfig

	Blockchains BlockchainConfig

	RedisConfig RedisConfig
}

type DBConfig struct {
	Host     string `envconfig:"PG_HOST" required:"true"`
	Port     string `envconfig:"PG_PORT" required:"true"`
	User     string `envconfig:"PG_USER" required:"true"`
	Password string `envconfig:"PG_PASSWORD" required:"true"`
}

type BlockchainConfig struct {
	PolygonAddr string `envconfig:"POLYGON_ADDR" required:"true"`
	BscAddr     string `envconfig:"BSC_ADDR" required:"true"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" required:"true"`
	Port     string `envconfig:"REDIS_PORT" required:"true"`
	Password string `envconfig:"REDIS_PASSWORD" required:"true"`
}

func (cfg *Config) GetRedisDSN() string {
	return fmt.Sprintf("redis://:%s@%s:%s", cfg.RedisConfig.Password, cfg.RedisConfig.Host, cfg.RedisConfig.Port)
}

func (cfg *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.User)
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
