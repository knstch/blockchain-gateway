package blockchain

import "math/big"

type WalletWithBalance struct {
	NativeBalance string
	Tokens        []Token
}

type Token struct {
	Balance string
	Symbol  string
}
type TokenInfo struct {
	Balance *big.Float
	Symbol  string
}
