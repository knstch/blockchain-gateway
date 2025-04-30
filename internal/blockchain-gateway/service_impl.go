package blockchain

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/knstch/subtrack-libs/svcerrs"
	"math/big"
	"strings"

	"github.com/knstch/subtrack-libs/log"

	"wallets-service/config"
	"wallets-service/internal/domain/enum"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrNoContractFound = fmt.Errorf("no contract found")
	ErrUnknownNetwork  = fmt.Errorf("unknown network: %w", svcerrs.ErrDataNotFound)
	ErrCantGetBaseFee  = fmt.Errorf("can't get base fee")
)

const erc20ABIJSON = `[
  {"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}
]`

type ClientImpl struct {
	Polygon  Chain
	Bsc      Chain
	erc20ABI abi.ABI
	lg       *log.Logger
	redis    *redis.Client
}

type Chain struct {
	Client  *ethclient.Client
	ChainID *big.Int
}

func NewClient(cfg *config.Config, logger *log.Logger, redis *redis.Client) (*ClientImpl, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		return nil, err
	}

	polygonClient, err := ethclient.Dial(cfg.Blockchains.PolygonAddr)
	if err != nil {
		return nil, fmt.Errorf("ethclient.Dial: %w", err)
	}
	polygonChainID, err := polygonClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("polygonClient.ChainID: %w", err)
	}

	bscClient, err := ethclient.Dial(cfg.Blockchains.BscAddr)
	if err != nil {
		return nil, fmt.Errorf("ethclient.Dial: %w", err)
	}
	bscChainID, err := bscClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("bscClient.ChainID: %w", err)
	}

	return &ClientImpl{
		Polygon: Chain{
			Client:  polygonClient,
			ChainID: polygonChainID,
		},
		Bsc: Chain{
			Client:  bscClient,
			ChainID: bscChainID,
		},
		erc20ABI: erc20ABI,
		lg:       logger,
		redis:    redis,
	}, nil
}

func (c *ClientImpl) getClient(network enum.Network) *ethclient.Client {
	switch network {
	case enum.PolygonNetwork:
		return c.Polygon.Client
	case enum.BscNetwork:
		return c.Bsc.Client
	default:
		return nil
	}
}

func isNoContractCodeError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no contract code at given address")
}

func (c *ClientImpl) getChainID(network enum.Network) *big.Int {
	switch network {
	case enum.PolygonNetwork:
		return c.Polygon.ChainID
	case enum.BscNetwork:
		return c.Bsc.ChainID
	default:
		return nil
	}
}

func (c *ClientImpl) getBaseGasFee(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("client.HeaderByNumber: %w", err)
	}

	if header.BaseFee == nil {
		return nil, ErrCantGetBaseFee
	}

	return header.BaseFee, nil
}
