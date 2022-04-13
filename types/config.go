package types

import (
	"io/ioutil"
	"path"
	"strings"

	junocmd "github.com/forbole/juno/v3/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configFileName = "config.yaml"
)

type Config struct {
	Report   *ReportConfig    `yaml:"report"`
	Chains   []*ChainConfig   `yaml:"chains"`
	Accounts []*AccountConfig `yaml:"accounts"`
}

func (c *Config) GetChainsList() []string {
	var chains []string
	for _, chain := range c.Chains {
		chains = append(chains, chain.Name)
	}
	return chains
}

func (c *Config) GetChainConfig(chainName string) *ChainConfig {
	for _, chain := range c.Chains {
		if strings.EqualFold(chain.Name, chainName) {
			return chain
		}
	}
	return nil
}

func (c *Config) GetChainAddresses(chainName string) (addresses []string, found bool) {
	for _, account := range c.Accounts {
		if strings.EqualFold(account.Chain, chainName) {
			return account.Addresses, true
		}
	}
	return nil, false
}

type ReportConfig struct {
	Currency string `yaml:"currency"`
}

type ChainConfig struct {
	Name               string `yaml:"name"`
	RPCAddress         string `yaml:"rpc_address"`
	GRPCAddress        string `yaml:"grpc_address"`
	AuthorizationToken string `yaml:"authorization,omitempty"`
}

type AccountConfig struct {
	Chain     string   `yaml:"chain"`
	Addresses []string `yaml:"addresses"`
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
