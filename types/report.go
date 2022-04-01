package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChainReport struct {
	ChainName    string         `yaml:"chain" json:"chain"`
	FirstBalance *BalanceReport `yaml:"first_balance" json:"first_balance"`
	LastBalance  *BalanceReport `yaml:"last_balance" json:"last_balance"`
}

func NewChaiReport(chainName string, firstBalance *BalanceReport, lastBalance *BalanceReport) *ChainReport {
	return &ChainReport{
		ChainName:    chainName,
		FirstBalance: firstBalance,
		LastBalance:  lastBalance,
	}
}

type BalanceReport struct {
	Timestamp time.Time `yaml:"timestamp" json:"timestamp"`
	Address   string    `yaml:"address" json:"address"`
	Amount    []Amount  `yaml:"amount" json:"amount"`
}

func NewBalanceReport(timestamp time.Time, address string, amount []Amount) *BalanceReport {
	return &BalanceReport{
		Timestamp: timestamp,
		Address:   address,
		Amount:    amount,
	}
}

type Amount struct {
	Coin  sdk.Coin `yaml:"coin" json:"coin"`
	Value string   `yaml:"value" json:"value"`
}

func NewAmount(coin sdk.Coin, value string) Amount {
	return Amount{
		Coin:  coin,
		Value: value,
	}
}
