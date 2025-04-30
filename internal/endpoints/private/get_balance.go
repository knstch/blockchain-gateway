package public

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/subtrack-libs/svcerrs"
	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets"

	public "github.com/knstch/blockchain-gateway-api/public"
)

func MakeGetBalanceEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return c.GetBalance(ctx, request.(*public.GetBalanceRequest))
	}
}

func (c *Controller) GetBalance(ctx context.Context, req *public.GetBalanceRequest) (*public.GetBalanceResponse, error) {
	network := enum.ConvertPublicNetworkToService(req.Network)
	if network == enum.UnknownNetwork {
		return nil, fmt.Errorf("unknown network: %w", svcerrs.ErrDataNotFound)
	}

	balance, err := c.svc.GetBalance(ctx, network)
	if err != nil {
		return nil, fmt.Errorf("svc.GetBalance: %w", err)
	}

	return &public.GetBalanceResponse{
		NativeBalance: balance.NativeBalance,
		Tokens:        convertServiceTokenBalanceToTransport(balance.Tokens),
	}, nil
}

func convertServiceTokenBalanceToTransport(wallet []wallets.Token) []*public.Token {
	tokens := make([]*public.Token, 0, len(wallet))

	for _, token := range wallet {
		tokens = append(tokens, &public.Token{
			Balance: token.Balance,
			Symbol:  token.Symbol,
		})
	}

	return tokens
}
