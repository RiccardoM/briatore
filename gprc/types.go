package gprc

import (
	"encoding/hex"
	"strings"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type StatusRequest struct{}

type NodeInfo struct {
	Network string `json:"network"`
}

type SyncInfo struct {
	LatestBlockHeight int64     `json:"latest_block_height,string"`
	LatestBlockTime   time.Time `json:"latest_block_time"`
}

type StatusResult struct {
	NodeInfo *NodeInfo `json:"node_info"`
	SyncInfo *SyncInfo `json:"sync_info"`
}

type BlockRequest struct {
	Height *int64 `json:"height,string,omitempty"`
}

type BlockHeader struct {
	ChainID string    `json:"chain_id"`
	Height  int64     `json:"height,string"`
	Time    time.Time `json:"time"`
}

type Block struct {
	BlockHeader `json:"header"`
}

type BlockResult struct {
	Block *Block `json:"block"`
}

type BlockResultsRequest struct {
	Height *int64 `json:"height,string,omitempty"`
}

type ResponseDeliverTx struct {
	Code      uint32            `json:"code"`
	Data      []byte            `json:"data"`
	Log       string            `json:"log"`
	GasWanted int64             `json:"gas_wanted,string"`
	GasUsed   int64             `json:"gas_used,string"`
	Events    []abcitypes.Event `json:"events"`
}

func (resp ResponseDeliverTx) IsOK() bool {
	return resp.Code == 0
}

type BlockResultsResult struct {
	Height           int64               `json:"height,string"`
	TxsResults       []ResponseDeliverTx `json:"txs_results"`
	BeginBlockEvents []abcitypes.Event   `json:"begin_block_events"`
	EndBlockEvents   []abcitypes.Event   `json:"end_block_events"`
}

type ABCIQueryRequest struct {
	Path   string         `json:"path"`
	Data   bytes.HexBytes `json:"data"`
	Height int64          `json:"height,string"`
	Prove  bool           `json:"prove"`
}

type ABCIQueryResponse struct {
	Code   uint32 `json:"code"`
	Log    string `json:"log"`
	Key    []byte `json:"key"`
	Value  []byte `json:"value"`
	Height int64  `json:"height,string"`
}

func (resp ABCIQueryResponse) IsOK() bool {
	return resp.Code == 0
}

type ABCIQueryResult struct {
	Response ABCIQueryResponse `json:"response"`
}

type BroadcastTxCommitRequest struct {
	Tx []byte `json:"tx"`
}

type ResponseCheckTx struct {
	Code      uint32            `json:"code"`
	Data      []byte            `json:"data"`
	Log       string            `json:"log"`
	GasWanted int64             `json:"gas_wanted,string"`
	GasUsed   int64             `json:"gas_used,string"`
	Events    []abcitypes.Event `json:"events"`
}

func (resp ResponseCheckTx) IsOK() bool {
	return resp.Code == 0
}

type BroadcastTxCommitResult struct {
	CheckTx   ResponseCheckTx   `json:"check_tx"`
	DeliverTx ResponseDeliverTx `json:"deliver_tx"`
	Hash      bytes.HexBytes    `json:"hash"`
	Height    int64             `json:"height,string"`
}

// copied from sdk
func NewResponseFormatBroadcastTxCommit(res *BroadcastTxCommitResult) *sdk.TxResponse {
	if res == nil {
		return nil
	}
	if !res.CheckTx.IsOK() {
		return newTxResponseCheckTx(res)
	}
	return newTxResponseDeliverTx(res)
}

// newTxResponseCheckTx has been copied from the Cosmos SDK
func newTxResponseCheckTx(res *BroadcastTxCommitResult) *sdk.TxResponse {
	if res == nil {
		return nil
	}
	var txHash string
	if res.Hash != nil {
		txHash = strings.ToUpper(hex.EncodeToString(res.Hash))
	}
	parsedLogs, _ := sdk.ParseABCILogs(res.CheckTx.Log)
	return &sdk.TxResponse{
		Height:    res.Height,
		TxHash:    txHash,
		Code:      res.CheckTx.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.CheckTx.Data)),
		RawLog:    res.CheckTx.Log,
		Logs:      parsedLogs,
		GasWanted: res.CheckTx.GasWanted,
		GasUsed:   res.CheckTx.GasUsed,
		Events:    res.CheckTx.Events,
	}
}

// newTxResponseDeliverTx has been copied from the Cosmos SDK
func newTxResponseDeliverTx(res *BroadcastTxCommitResult) *sdk.TxResponse {
	if res == nil {
		return nil
	}
	var txHash string
	if res.Hash != nil {
		txHash = strings.ToUpper(hex.EncodeToString(res.Hash))
	}
	parsedLogs, _ := sdk.ParseABCILogs(res.DeliverTx.Log)
	return &sdk.TxResponse{
		Height:    res.Height,
		TxHash:    txHash,
		Code:      res.DeliverTx.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.DeliverTx.Data)),
		RawLog:    res.DeliverTx.Log,
		Logs:      parsedLogs,
		GasWanted: res.DeliverTx.GasWanted,
		GasUsed:   res.DeliverTx.GasUsed,
		Events:    res.DeliverTx.Events,
	}
}
