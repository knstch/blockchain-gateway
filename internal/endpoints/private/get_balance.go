package private

import (
	"blockchain-gateway/internal/blockchain-gateway"
	"blockchain-gateway/internal/helpers"
	"context"
	"fmt"
	"github.com/knstch/subtrack-libs/enum"
	"github.com/knstch/subtrack-libs/svcerrs"
	"github.com/knstch/subtrack-libs/tracing"

	private "github.com/knstch/blockchain-gateway-api/private"
)

func (c *Controller) GetBalance(ctx context.Context, req *private.GetBalanceRequest) (*private.GetBalanceResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "private: GetBalance")
	defer span.End()

	network := helpers.ConvertPublicNetworkToService(req.Network)
	if network == enum.UnknownNetwork {
		return nil, fmt.Errorf("unknown network: %w", svcerrs.ErrDataNotFound)
	}

	balance, err := c.svc.GetBalance(ctx, network, req.PublicAddress, req.TokenAddresses)
	if err != nil {
		return nil, fmt.Errorf("svc.GetBalance: %w", err)
	}

	return &private.GetBalanceResponse{
		NativeBalance: balance.NativeBalance,
		Tokens:        convertServiceTokenBalanceToTransport(balance.Tokens),
	}, nil
}

func convertServiceTokenBalanceToTransport(wallet []blockchain.Token) []*private.Token {
	tokens := make([]*private.Token, 0, len(wallet))

	for _, token := range wallet {
		tokens = append(tokens, &private.Token{
			Balance: token.Balance,
			Symbol:  token.Symbol,
		})
	}

	return tokens
}
