package public

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	public "github.com/knstch/blockchain-gateway-api/public"
	"github.com/knstch/subtrack-libs/log"
	"github.com/knstch/subtrack-libs/middleware"
	"github.com/knstch/subtrack-libs/transport"
	"net/http"

	"wallets-service/config"
	"wallets-service/internal/wallets"

	"github.com/knstch/subtrack-libs/endpoints"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
}

type Controller struct {
	svc wallets.Wallets
	lg  *log.Logger
	cfg *config.Config
}

func NewController(svc wallets.Wallets, lg *log.Logger, cfg *config.Config) *Controller {
	return &Controller{
		svc: svc,
		cfg: cfg,
		lg:  lg,
	}
}

func (c *Controller) Endpoints() []endpoints.Endpoint {
	mdw := []middleware.Middleware{middleware.WithCookieAuth(c.cfg.JwtSecret)}

	return []endpoints.Endpoint{
		{
			Method:  http.MethodPost,
			Path:    "/getBalance",
			Handler: MakeGetBalanceEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.GetBalanceRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Mdw:     mdw,
		},
	}
}
