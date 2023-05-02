package main

type TickerMessage struct {
	Data Ticker `json:"data"`
}

type Ticker struct {
	Sym                string      `json:"s"`
	Last               interface{} `json:"c"`
	Bid                string      `json:"b"`
	Ask                string      `json:"a"`
	Volume             string      `json:"q"`
	PriceChangePercent string      `json:"P"`
}

type Symbol struct {
	Symbol string `json:"symbol"`
}
