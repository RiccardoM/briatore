package types

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const (
	blocksDataFile = "blocks.json"
)

type BlocksData map[string]BlockData

type BlockData struct {
	Height    int64     `json:"height"`
	Timestamp time.Time `json:"timestamp"`
}

func NewBlockData(height int64, timestamp time.Time) BlockData {
	return BlockData{
		Timestamp: timestamp,
		Height:    height,
	}
}

func (b BlockData) IsZero() bool {
	return b.Height == 0
}

func GetBlocksData() (BlocksData, error) {
	bz, err := ioutil.ReadFile(path.Join(homePath, blocksDataFile))
	if os.IsNotExist(err) {
		return BlocksData{}, nil
	}
	if err != nil {
		return nil, err
	}

	var data BlocksData
	return data, json.Unmarshal(bz, &data)
}

func WriteBlocksData(data BlocksData) error {
	bz, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(homePath, blocksDataFile), bz, 0600)
}
