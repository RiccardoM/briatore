package cosmos

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	tmjson "github.com/cometbft/cometbft/libs/json"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	tmtypes "github.com/cometbft/cometbft/types"
)

type Client struct {
	ctx    context.Context
	client *httpclient.HTTP
}

func NewClient(rpcAddress string) (*Client, error) {
	httpClient, err := jsonrpcclient.DefaultHTTPClient(rpcAddress)
	if err != nil {
		return nil, err
	}

	// Tweak the transport
	httpTransport, ok := (httpClient.Transport).(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("invalid HTTP Transport: %T", httpTransport)
	}
	httpTransport.MaxConnsPerHost = 20

	rpcClient, err := httpclient.NewWithClient(rpcAddress, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	err = rpcClient.Start()
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:    context.Background(),
		client: rpcClient,
	}, nil
}

// Genesis returns the genesis file
func (cp *Client) Genesis() (*tmctypes.ResultGenesis, error) {
	res, err := cp.client.Genesis(cp.ctx)
	if err != nil && strings.Contains(err.Error(), "use the genesis_chunked API instead") {
		return cp.getGenesisChunked()
	}
	return res, err
}

// getGenesisChunked gets the genesis data using the chinked API instead
func (cp *Client) getGenesisChunked() (*tmctypes.ResultGenesis, error) {
	bz, err := cp.getGenesisChunksStartingFrom(0)
	if err != nil {
		return nil, err
	}

	var genDoc *tmtypes.GenesisDoc
	err = tmjson.Unmarshal(bz, &genDoc)
	if err != nil {
		return nil, err
	}

	return &tmctypes.ResultGenesis{Genesis: genDoc}, nil
}

// getGenesisChunksStartingFrom returns all the genesis chunks data starting from the chunk with the given id
func (cp *Client) getGenesisChunksStartingFrom(id uint) ([]byte, error) {
	res, err := cp.client.GenesisChunked(cp.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting genesis chunk %d", id)
	}

	bz, err := base64.StdEncoding.DecodeString(res.Data)
	if err != nil {
		return nil, fmt.Errorf("error while decoding genesis chunk %d out of %d", id, res.TotalChunks)
	}

	if id == uint(res.TotalChunks-1) {
		return bz, nil
	}

	nextChunk, err := cp.getGenesisChunksStartingFrom(id + 1)
	if err != nil {
		return nil, err
	}

	return append(bz, nextChunk...), nil
}

// --------------------------------------------------------------------------------------------------------------------

// Block returns the block having the given height
func (cp *Client) Block(height int64) (*tmctypes.ResultBlock, error) {
	return cp.client.Block(cp.ctx, &height)
}

// --------------------------------------------------------------------------------------------------------------------

// LatestHeight implements node.Node
func (cp *Client) LatestHeight() (int64, error) {
	status, err := cp.client.Status(cp.ctx)
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}
