package types

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	configFileName = "config.yaml"
)

type Config struct {
	Report *ReportConfig  `yaml:"report"`
	Chains []*ChainConfig `yaml:"chains"`
}

type ReportConfig struct {
	Currency string `yaml:"currency"`
}

type ChainConfig struct {
	Name           string `yaml:"name"`
	RPCAddress     string `yaml:"rpcAddress"`
	AssetName      string `yaml:"asset"`
	Bech32Prefix   string `yaml:"bech32Prefix"`
	MinBlockHeight int64  `yaml:"minBlockHeight"`
}

type AccountConfig struct {
	Chain     string   `yaml:"chain"`
	Addresses []string `yaml:"addresses"`
}

// ReadConfig reads the config from the given command
func ReadConfig(cmd *cobra.Command) (*Config, error) {
	home, err := cmd.Flags().GetString("home")
	if err != nil {
		return nil, err
	}
	HomePath = home

	cfgPath := path.Join(HomePath, configFileName)
	log.Debug().Str("home", cfgPath).Msg("reading config file")

	bz, err := os.ReadFile(cfgPath)
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
