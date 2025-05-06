package enum

type Op string

const (
	AddWallet    Op = "addWallet"
	RemoveWallet Op = "removeWallet"
)

func (op Op) String() string {
	return string(op)
}
