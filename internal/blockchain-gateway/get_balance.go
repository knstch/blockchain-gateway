package blockchain

import (
	"blockchain-gateway/internal/abi"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knstch/subtrack-libs/enum"
	"github.com/knstch/subtrack-libs/tracing"
	"math/big"
)

func (svc *ServiceImpl) getNativeBalance(ctx context.Context, walletAddr string, network enum.Network) (*big.Float, error) {
	ctx, span := tracing.StartSpan(ctx, "service: getNativeBalance")
	defer span.End()

	client := svc.getClient(network)
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

func (svc *ServiceImpl) getTokenBalanceAndInfo(ctx context.Context, walletAddr, tokenAddr string, network enum.Network) (TokenInfo, error) {
	ctx, span := tracing.StartSpan(ctx, "service: getTokenBalanceAndInfo")
	defer span.End()

	client := svc.getClient(network)
	if client == nil {
		return TokenInfo{}, ErrUnknownNetwork
	}

	walletAddressCommon := common.HexToAddress(walletAddr)
	tokenAddressCommon := common.HexToAddress(tokenAddr)

	contract := bind.NewBoundContract(tokenAddressCommon, abi.GetErc20Abi(), client, client, client)

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
	ctx, span := tracing.StartSpan(ctx, "service: getBalance")
	defer span.End()

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
	ctx, span := tracing.StartSpan(ctx, "service: getSymbol")
	defer span.End()

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

func (svc *ServiceImpl) GetBalance(ctx context.Context, network enum.Network, publicAddr string, tokenAddresses []string) (WalletWithBalance, error) {
	ctx, span := tracing.StartSpan(ctx, "service: GetBalance")
	defer span.End()

	nativeBalance, err := svc.getNativeBalance(ctx, publicAddr, network)
	if err != nil {
		return WalletWithBalance{}, fmt.Errorf("blockchain.GetNativeBalance: %w", err)
	}

	tokens := make([]Token, 0, len(tokenAddresses))
	for _, address := range tokenAddresses {
		token, err := svc.getTokenBalanceAndInfo(ctx, publicAddr, address, network)
		if err != nil {
			if errors.Is(err, ErrNoContractFound) {
				continue
			}
			return WalletWithBalance{}, fmt.Errorf("blockchain.GetTokenBalanceAndInfo: %w", err)
		}

		tokens = append(tokens, Token{
			Symbol:  token.Symbol,
			Balance: token.Balance.String(),
		})
	}

	balance := &WalletWithBalance{
		NativeBalance: nativeBalance.String(),
		Tokens:        tokens,
	}

	return *balance, nil
}
