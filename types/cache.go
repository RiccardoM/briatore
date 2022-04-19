package types

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

const (
	cacheFileName = "cache.json"
)

// Cache contains a list of [ChainName -> []CacheEntry] entries
type Cache struct {
	Blocks []BlockData `json:"blocks"`
	Prices []PriceData `json:"prices"`
}

func readCache() (Cache, error) {
	// Read the cache
	bz, err := ioutil.ReadFile(path.Join(HomePath, cacheFileName))
	if os.IsNotExist(err) {
		return Cache{}, nil
	}
	if err != nil {
		panic(err)
	}

	var data Cache
	return data, json.Unmarshal(bz, &data)
}

func writeCache(cache Cache) error {
	bz, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(HomePath, cacheFileName), bz, 0600)
}

// --------------------------------------------------------------------------------------------------------------------

type BlockData struct {
	ChainName string    `json:"chain"`
	Height    int64     `json:"height"`
	Timestamp time.Time `json:"timestamp"`
}

func NewBlockData(chainName string, height int64, timestamp time.Time) BlockData {
	return BlockData{
		ChainName: chainName,
		Timestamp: timestamp,
		Height:    height,
	}
}

func (b BlockData) IsZero() bool {
	return b.Height == 0
}

func GetBlockData(chainName string, timestamp time.Time) (data BlockData, found bool, err error) {
	cache, err := readCache()
	if err != nil {
		return
	}

	for _, block := range cache.Blocks {
		if strings.EqualFold(block.ChainName, chainName) && IsSameDay(block.Timestamp, timestamp) {
			return block, true, nil
		}
	}

	return BlockData{}, false, nil
}

func CacheBlockData(data BlockData) error {
	cache, err := readCache()
	if err != nil {
		return err
	}

	cache.Blocks = append(cache.Blocks, data)
	return writeCache(cache)
}

// --------------------------------------------------------------------------------------------------------------------

type PriceData struct {
	CoinGeckoID string    `json:"coinGeckoID"`
	Price       float64   `json:"price"`
	Timestamp   time.Time `json:"timestamp"`
	Currency    string    `json:"currency"`
}

func NewPriceData(coinGeckoID string, price float64, currency string, timestamp time.Time) PriceData {
	return PriceData{
		CoinGeckoID: coinGeckoID,
		Price:       price,
		Currency:    currency,
		Timestamp:   timestamp,
	}
}

func GetPriceData(coinGeckoID string, currency string, timestamp time.Time) (data PriceData, found bool, err error) {
	cache, err := readCache()
	if err != nil {
		return
	}

	for _, price := range cache.Prices {
		if price.CoinGeckoID == coinGeckoID && price.Currency == currency && IsSameDay(price.Timestamp, timestamp) {
			return price, true, nil
		}
	}

	return PriceData{}, false, nil
}

func CachePriceData(data PriceData) error {
	cache, err := readCache()
	if err != nil {
		return err
	}

	cache.Prices = append(cache.Prices, data)
	return writeCache(cache)
}

// --------------------------------------------------------------------------------------------------------------------

func IsSameDay(first, second time.Time) bool {
	return first.Year() == second.Year() &&
		first.Month() == second.Month() &&
		first.Day() == second.Day()
}
