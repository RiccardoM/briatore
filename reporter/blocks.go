package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/tendermint/tendermint/rpc/jsonrpc/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
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

	latestBlock, err := r.getBlockOrLatestHeight(latestHeight)
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
	log.Trace().Int64("min height", minHeight).Int64("max height", maxHeight).Time("timestamp", timestamp).Msg("binary search")

	minBlock, err := r.getBlockOrMinHeight(minHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting min block: %s", err)
	}

	if minBlock.Block.Time.Equal(timestamp) {
		// The genesis has the same timestamp of what we are searching for
		return minBlock.Block, nil
	}

	maxBlock, err := r.getBlockOrMinHeight(maxHeight)
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

	avgBlock, err := r.getBlockOrMinHeight(avgHeight)
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

// getBlockOrMinHeight gets the block at the given height, or the min height available if not found
func (r *Reporter) getBlockOrMinHeight(height int64) (*tmctypes.ResultBlock, error) {
	block, err := r.node.Block(height)

	if err != nil {
		if rpcErr, ok := err.(*types.RPCError); ok && strings.Contains(rpcErr.Data, "lowest height") {
			var lowestHeight int64
			_, err = fmt.Sscanf(rpcErr.Data, "height %d is not available, lowest height is %d", &height, &lowestHeight)
			if err != nil {
				return nil, err
			}

			log.Debug().Str("chain", r.cfg.Name).
				Int64("height", height).Int64("lowest height", lowestHeight).
				Msg("height not found, getting lowest height")

			return r.node.Block(lowestHeight)
		}

		return nil, err
	}

	return block, nil
}

// getBlockOrLatestHeight gets the block at the given height, or the max height available if not found
func (r *Reporter) getBlockOrLatestHeight(height int64) (*tmctypes.ResultBlock, error) {
	block, err := r.node.Block(height)

	if err != nil {
		if rpcErr, ok := err.(*types.RPCError); ok && strings.Contains(rpcErr.Data, "current blockchain height") {
			var maxHeight int64
			_, err = fmt.Sscanf(rpcErr.Data, "height %d must be less than or equal to the current blockchain height %d", &height, &maxHeight)
			if err != nil {
				return nil, err
			}

			log.Debug().Str("chain", r.cfg.Name).
				Int64("height", height).Int64("max height", maxHeight).
				Msg("height not found, getting max height")

			return r.node.Block(maxHeight)
		}

		return nil, err
	}

	return block, nil
}
