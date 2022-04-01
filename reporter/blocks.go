package reporter

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// getBlockNearTimestamp returns the block nearest the given timestamp.
// To do this we use the binary search between the genesis height and the latest block time.
func (r *Reporter) getBlockNearTimestamp(timestamp time.Time) (*tmtypes.Block, error) {
	log.Debug().Str("chain", r.cfg.Name).Time("timestamp", timestamp).Msg("getting block near timestamp")

	genesis, err := r.node.Genesis()
	if err != nil {
		return nil, fmt.Errorf("error while getting the genesis: %s", err)
	}

	if timestamp.Before(genesis.Genesis.GenesisTime) {
		return nil, nil
	}

	genesisHeight := genesis.Genesis.InitialHeight

	latestHeight, err := r.node.LatestHeight()
	if err != nil {
		return nil, fmt.Errorf("error while getting latest height: %s", err)
	}

	latestBlock, err := r.node.Block(latestHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting latest block: %s", err)
	}

	if timestamp.After(latestBlock.Block.Time) {
		return nil, fmt.Errorf("%s is after latest block time", timestamp)
	}

	// Perform the binary search
	block, err := r.binarySearchBlock(genesisHeight, latestHeight, timestamp)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("chain", r.cfg.Name).Time("timestamp", timestamp).Msgf("found block near timestamp: %d", block.Height)

	return block, nil
}

// binarySearchBlock performs a binary search between the given min and max heights,
// searching for the block that is closer to the given timestamp
func (r *Reporter) binarySearchBlock(minHeight, maxHeight int64, timestamp time.Time) (*tmtypes.Block, error) {
	log.Trace().Int64("ming height", minHeight).Int64("max height", maxHeight).Time("timestamp", timestamp).Msg("binary search")

	minBlock, err := r.node.Block(minHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting min block: %s", err)
	}

	if minBlock.Block.Time.Equal(timestamp) {
		// The genesis has the same timestamp of what we are searching for
		return minBlock.Block, nil
	}

	maxBlock, err := r.node.Block(maxHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting max block")
	}

	if maxBlock.Block.Time.Equal(timestamp) {
		// The latest block has the same timestamp we are searching for
		return maxBlock.Block, nil
	}

	if maxBlock.Block.Height-minBlock.Block.Height == 1 {
		// We've reached the point where we only have two blocks.
		// Now we need to find the one that is closer to the searched timestamp
		minDiff := timestamp.Sub(minBlock.Block.Time)
		maxDiff := maxBlock.Block.Time.Sub(timestamp)

		if minDiff < maxDiff {
			// The min block is closer to the timestamp than the max block
			return minBlock.Block, nil
		}

		// The max block is closer to the timestamp than the min block
		return maxBlock.Block, nil
	}

	avgHeight := (maxHeight + minHeight) / 2

	avgBlock, err := r.node.Block(avgHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting average block: %s", err)
	}

	if avgBlock.Block.Time.Equal(timestamp) {
		// The average block has the same timestamp as the one searched for
		return avgBlock.Block, nil
	}

	if avgBlock.Block.Time.After(timestamp) {
		// If the average block has the timestamp after the searched value, it means the searched
		// value is in between the min value and the average one
		maxHeight = avgBlock.Block.Height
	}

	if avgBlock.Block.Time.Before(timestamp) {
		// If the average block has the timestamp before the searched value, it means the searched
		// value is in between the average block and the max height
		minHeight = avgBlock.Block.Height
	}

	return r.binarySearchBlock(minHeight, maxHeight, timestamp)
}
