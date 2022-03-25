package types

import (
	"io/ioutil"
	"path"
	"strings"

	junocmd "github.com/forbole/juno/v3/cmd"
	"github.com/forbole/juno/v3/node/remote"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configFileName = "config.yaml"
)

type Config struct {
	Chains []*ChainConfig `yaml:"chains"`
	Report *ReportConfig  `yaml:"report"`
}

func (c *Config) GetChainConfig(chainName string) (*ChainConfig, bool) {
	for _, chain := range c.Chains {
		if strings.EqualFold(chain.Name, chainName) {
			return chain, true
		}
	}
	return nil, false
}

func (c *Config) GetChainsList() []string {
	names := make([]string, len(c.Chains))
	for i, chain := range c.Chains {
		names[i] = chain.Name
	}
	return names
}

type ChainConfig struct {
	Name    string         `yaml:"name"`
	Address string         `yaml:"address"`
	Node    remote.Details `yaml:"node"`
}
type ReportConfig struct {
	Currency string `yaml:"currency"`
	Coins    []Coin `yaml:"coins"`
}

func (c *ReportConfig) GetCoinGeckoID(denom string) (string, bool) {
	for _, coin := range c.Coins {
		if strings.EqualFold(coin.Denom, denom) {
			return coin.CoinGeckoID, true
		}
	}
	return "", false
}

type Coin struct {
	Denom       string `yaml:"denom"`
	CoinGeckoID string `yaml:"coingecko_id"`
}

// ReadConfig reads the config from the given command
func ReadConfig(cmd *cobra.Command) (*Config, error) {
	home, err := cmd.Flags().GetString(junocmd.FlagHome)
	if err != nil {
		return nil, err
	}

	cfgPath := path.Join(home, configFileName)
	log.Debug().Str("home", cfgPath).Msg("reading config file")

	bz, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(bz, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
