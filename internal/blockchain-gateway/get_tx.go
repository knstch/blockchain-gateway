package blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"wallets-service/internal/domain/enum"
)

func (c *ClientImpl) GetTransaction(ctx context.Context, txID string, network enum.Network) (*types.Transaction, error) {
	client := c.getClient(network)
	if client == nil {
		return nil, ErrUnknownNetwork
	}

	txHash := common.HexToHash(txID)
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("client.TransactionByHash: %w", err)
	}

	return tx, nil
}
