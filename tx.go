package walletwatch

type CoinNetwork int

const (
	BitcoinMainNet CoinNetwork = iota
)

type Transaction interface {
	// The blockchain network this transaction took place on
	Network() CoinNetwork

	// The network-unique identifier for this transaction
	Hash() string

	// The outputs involved in the transaction
	Outputs() map[string]int
}
