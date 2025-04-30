package blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"wallets-service/internal/domain/enum"
)

func (c *ClientImpl) GetNativeBalance(ctx context.Context, walletAddr string, network enum.Network) (*big.Float, error) {
	client := c.getClient(network)
	if client == nil {
		return nil, ErrUnknownNetwork
	}

	account := common.HexToAddress(walletAddr)

	balanceWei, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("client.BalanceAt: %w", err)
	}

	balanceCurrency := new(big.Float)
	balanceCurrency.SetString(balanceWei.String())
	ethValue := new(big.Float).Quo(balanceCurrency, big.NewFloat(1e18))

	return ethValue, nil
}

func (c *ClientImpl) GetTokenBalanceAndInfo(ctx context.Context, walletAddr, tokenAddr string, network enum.Network) (TokenInfo, error) {
	client := c.getClient(network)
	if client == nil {
		return TokenInfo{}, ErrUnknownNetwork
	}

	walletAddressCommon := common.HexToAddress(walletAddr)
	tokenAddressCommon := common.HexToAddress(tokenAddr)

	contract := bind.NewBoundContract(tokenAddressCommon, c.erc20ABI, client, client, client)

	balance, err := getBalance(ctx, walletAddressCommon, contract)
	if err != nil {
		return TokenInfo{}, err
	}

	symbol, err := getSymbol(ctx, contract)
	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Balance: balance,
		Symbol:  symbol,
	}, nil
}

func getBalance(ctx context.Context, walletAddr common.Address, contract *bind.BoundContract) (*big.Float, error) {
	var outBalance []interface{}
	err := contract.Call(&bind.CallOpts{Context: ctx}, &outBalance, "balanceOf", walletAddr)
	if err != nil {
		if isNoContractCodeError(err) {
			return nil, ErrNoContractFound
		}
		return nil, fmt.Errorf("contract.Call: %w", err)
	}

	balance, ok := outBalance[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected type in balanceOf result: %T", outBalance[0])
	}

	readableBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))

	return readableBalance, nil
}

func getSymbol(ctx context.Context, contract *bind.BoundContract) (string, error) {
	var outSymbol []interface{}
	err := contract.Call(&bind.CallOpts{Context: ctx}, &outSymbol, "symbol")
	if err != nil {
		if isNoContractCodeError(err) {
			return "", ErrNoContractFound
		}
		return "", fmt.Errorf("contract.Call: %w", err)
	}
	symbol, ok := outSymbol[0].(string)
	if !ok {
		symbol = ""
	}

	return symbol, nil
}
