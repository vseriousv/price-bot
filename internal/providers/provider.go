package providers

import "errors"

type Price string
type Ticker string

type IProvider interface {
	GetPriceByTicker(ticker string) *Price
	GetTickersList() []Ticker
}

func GetProvider(name string) (IProvider, error) {
	switch name {
	case "kucoin":
		p := &KucoinProvider{}
		p.SetParams()
		return p, nil
	default:
		return nil, errors.New("the provider is not found")
	}
}
