package providers

import (
	"encoding/json"
	"fmt"
	"github.com/vseriousv/price-bot/internal/models"
	"io"
	"log"
	"net/http"
)

type KucoinProvider struct {
	models.Provider
}

func (p *KucoinProvider) SetParams() {
	p.Name = "kucoin"
	p.ApiUrl = "https://api.kucoin.com"
}

func (p *KucoinProvider) GetPriceByTicker(ticker string) *Price {
	query := fmt.Sprintf(
		"%s/api/v1/market/orderbook/level1?symbol=%s",
		p.ApiUrl,
		ticker,
	)
	log.Printf("[provider/%s] :: %s", p.Name, query)

	resp, err := http.Get(query)
	if err != nil {
		log.Println(err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	tickerStruct := struct {
		Code string          `json:"code"`
		Data TickerMarketDto `json:"data"`
	}{}

	err = json.Unmarshal(body, &tickerStruct)
	if err != nil {
		log.Println(err)
		return nil
	}
	var price = Price(tickerStruct.Data.Price)
	return &price
}

func (p *KucoinProvider) GetTickersList() []Ticker {
	var arr []Ticker
	query := fmt.Sprintf(
		"%s/api/v1/market/allTickers",
		p.ApiUrl,
	)
	log.Printf("[provider/%s] :: %s", p.Name, query)

	resp, err := http.Get(query)
	if err != nil {
		log.Println(err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	allTickersStruct := struct {
		Code string `json:"code"`
		Data struct {
			Time   int64       `json:"time"`
			Ticker []TickerDto `json:"ticker"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &allTickersStruct)
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, item := range allTickersStruct.Data.Ticker {
		arr = append(arr, Ticker(item.Symbol))
	}

	log.Println("len(arr)", len(arr))
	return arr
}
