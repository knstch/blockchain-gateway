package private

import (
	"blockchain-gateway/config"
	"blockchain-gateway/internal/blockchain-gateway"
	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/subtrack-libs/endpoints"
	"github.com/knstch/subtrack-libs/log"

	privateApi "github.com/knstch/blockchain-gateway-api/private"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
}

type Controller struct {
	svc blockchain.Blockchain
	lg  *log.Logger
	cfg *config.Config

	privateApi.UnimplementedBlockchainGatewayServer
}

func NewController(svc blockchain.Blockchain, lg *log.Logger, cfg *config.Config) *Controller {
	return &Controller{
		svc: svc,
		cfg: cfg,
		lg:  lg,
	}
}

func (c *Controller) Endpoints() []endpoints.Endpoint {
	return []endpoints.Endpoint{}
}
