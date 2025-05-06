package abi

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

const erc20ABIJSON = `[
  {"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}
]`

func GetErc20Abi() abi.ABI {
	erc20ABI, _ := abi.JSON(strings.NewReader(erc20ABIJSON))
	
	return erc20ABI
}
