package cosmos

import (
	"context"
	"fmt"
	"net/http"

	httpclient "github.com/cometbft/cometbft/rpc/client/http"
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

// MinHeight returns the minimum height of the chain
func (cp *Client) MinHeight() (int64, error) {
	res, err := cp.client.Status(cp.ctx)
	if err != nil {
		return 0, err
	}
	return res.SyncInfo.EarliestBlockHeight, nil
}

// LatestHeight returns the latest height of the chain
func (cp *Client) LatestHeight() (int64, error) {
	status, err := cp.client.Status(cp.ctx)
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// Block returns the block having the given height
func (cp *Client) Block(height int64) (*tmtypes.Block, error) {
	block, err := cp.client.Block(cp.ctx, &height)
	if err != nil {
		return nil, err
	}
	return block.Block, nil
}
