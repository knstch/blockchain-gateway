package blockchain

import (
	"blockchain-gateway/config"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/knstch/subtrack-libs/enum"
	"math/big"

	"github.com/knstch/subtrack-libs/log"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Blockchain interface {
	GetBalance(ctx context.Context, network enum.Network, publicAddr string, tokenAddresses []string) (WalletWithBalance, error)
}

type ServiceImpl struct {
	Bsc   Chain
	lg    *log.Logger
	redis *redis.Client
}

type Chain struct {
	Client  *ethclient.Client
	ChainID *big.Int
}

func NewService(cfg *config.Config, logger *log.Logger) (*ServiceImpl, error) {
	bscClient, err := ethclient.Dial(cfg.Blockchains.BscAddr)
	if err != nil {
		return nil, fmt.Errorf("ethclient.Dial: %w", err)
	}
	bscChainID, err := bscClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("bscClient.ChainID: %w", err)
	}

	return &ServiceImpl{
		Bsc: Chain{
			Client:  bscClient,
			ChainID: bscChainID,
		},
		lg: logger,
	}, nil
}

func (svc *ServiceImpl) getClient(network enum.Network) *ethclient.Client {
	switch network {
	case enum.BscNetwork:
		return svc.Bsc.Client
	default:
		return nil
	}
}

func (svc *ServiceImpl) getChainID(network enum.Network) *big.Int {
	switch network {
	case enum.BscNetwork:
		return svc.Bsc.ChainID
	default:
		return nil
	}
}

func (svc *ServiceImpl) getBaseGasFee(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("client.HeaderByNumber: %w", err)
	}

	if header.BaseFee == nil {
		return nil, ErrCantGetBaseFee
	}

	return header.BaseFee, nil
}
