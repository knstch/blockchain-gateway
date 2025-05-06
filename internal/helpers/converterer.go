package helpers

import "github.com/knstch/subtrack-libs/enum"

func ConvertPublicNetworkToService(network string) enum.Network {
	switch network {
	case "bsc":
		return enum.BscNetwork
	default:
		return enum.UnknownNetwork
	}
}
