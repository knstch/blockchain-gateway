package enum

import (
	public "github.com/knstch/wallets-api/public"
)

type Network string

const (
	UnknownNetwork Network = "unknown"
	PolygonNetwork Network = "polygon"
	BscNetwork     Network = "bsc"
)

func ConvertPublicNetworkToService(network public.Network) Network {
	switch network {
	case public.Network_POLYGON:
		return PolygonNetwork
	case public.Network_BSC:
		return BscNetwork
	default:
		return UnknownNetwork
	}
}
