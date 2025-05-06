package enum

type Network string

func (n Network) String() string {
	return string(n)
}

const (
	UnknownNetwork Network = "unknown"
	PolygonNetwork Network = "polygon"
	BscNetwork     Network = "bsc"
)
