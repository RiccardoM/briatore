package types

import "fmt"

type HistoryResponse struct {
	MarketData *MarketData `json:"market_data"`
}

func (h HistoryResponse) GetCoinPrice(currency string) (float64, error) {
	if h.MarketData == nil {
		return 0, nil
	}

	price, ok := h.MarketData.CurrentPrice[currency]
	if !ok {
		return 0, fmt.Errorf("invalid currency: %s", currency)
	}

	return price, nil
}

type MarketData struct {
	CurrentPrice map[string]float64 `json:"current_price"`
}
