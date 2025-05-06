package blockchain

import (
	"fmt"
	"github.com/knstch/subtrack-libs/svcerrs"
	"strings"
)

var (
	ErrNoContractFound = fmt.Errorf("no contract found")
	ErrUnknownNetwork  = fmt.Errorf("unknown network: %w", svcerrs.ErrDataNotFound)
	ErrCantGetBaseFee  = fmt.Errorf("can't get base fee")
)

func isNoContractCodeError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no contract code at given address")
}
