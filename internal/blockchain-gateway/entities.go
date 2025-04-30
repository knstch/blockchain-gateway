package blockchain

import "math/big"

type Wallet struct {
	PublicAddr     string
	TokenAddresses []string
}

type TokenInfo struct {
	Balance *big.Float
	Symbol  string
}
