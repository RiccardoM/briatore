package reporter

import (
	tmtypes "github.com/cometbft/cometbft/types"
)

type CosmosClient interface {
	MinHeight() (int64, error)
	LatestHeight() (int64, error)
	Block(height int64) (*tmtypes.Block, error)
}
