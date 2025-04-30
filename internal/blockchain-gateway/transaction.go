package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"wallets-service/internal/domain/enum"
)

func (c *ClientImpl) BuildAndSendTx(ctx context.Context, network enum.Network, privateKey *ecdsa.PrivateKey,
	walletAddress common.Address, to common.Address, data []byte) (*types.Transaction, error) {
	client := c.getClient(network)
	if client == nil {
		return nil, ErrUnknownNetwork
	}
	chainID := c.getChainID(network)

	tx, err := c.buildTx(ctx, client, walletAddress, to, data, chainID)
	if err != nil {
		return nil, err
	}

	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, fmt.Errorf("SignTx: %w", err)
	}

	if err = client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("client.SendTransaction: %w", err)
	}

	return signedTx, nil
}

func (c *ClientImpl) buildTx(ctx context.Context,
	client *ethclient.Client, walletAddress common.Address, to common.Address, data []byte, chainID *big.Int) (*types.Transaction, error) {
	nonce, err := client.PendingNonceAt(ctx, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("client.PendingNonceAt: %w", err)
	}

	baseFee, err := c.getBaseGasFee(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("getBaseGasFee: %w", err)
	}
	priorityFee, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("client.SuggestGasPrice: %w", err)
	}

	maxFee := new(big.Int).Mul(baseFee, big.NewInt(2))
	maxFee = new(big.Int).Add(maxFee, priorityFee)

	multipliedGas := (baseFee.Uint64() / 100 * 10) + baseFee.Uint64()

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: priorityFee,
		GasFeeCap: maxFee,
		To:        &to,
		Data:      data,
		Gas:       multipliedGas,
	})

	return tx, nil
}
