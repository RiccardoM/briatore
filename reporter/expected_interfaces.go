package reporter

import (
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type CosmosClient interface {
	Genesis() (*tmctypes.ResultGenesis, error)
	LatestHeight() (int64, error)
	Block(height int64) (*tmctypes.ResultBlock, error)
}
