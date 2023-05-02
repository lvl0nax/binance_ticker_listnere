package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	redisDb := NewStorage("redis://localhost:6379/1")
	defer redisDb.Close()

	symbols := GetBinanceSymbols()
	strSymbols := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		strSymbols = append(strSymbols, strings.ToLower(symbol.Symbol))
	}
	chunkedSymbols := chunkSlice(strSymbols, 100)
	for _, chunkSymbols := range chunkedSymbols {
		go RunWebsocket(chunkSymbols, redisDb)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(8 * time.Minute)
	fmt.Println("Sleep Over.....")
	//fmt.Println(symbols, len(symbols))
}

func GetBinanceSymbols() []Symbol {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price")
	if err != nil {
		panic(err) // handle error
	}
	defer resp.Body.Close()

	symbols := make([]Symbol, 0, 2300)
	err = json.NewDecoder(resp.Body).Decode(&symbols)
	if err != nil {
		panic(err)
	}

	return symbols
}

func RunWebsocket(symbols []string, redisDb Storage) {
	streams := strings.Join(symbols, "@ticker/")
	wssUrl, err := url.Parse("wss://stream.binance.com/stream?streams=" + streams + "@ticker")
	if err != nil {
		panic(err)
	}

	conn, _, err := websocket.Dial(context.Background(), wssUrl.String(), nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close(websocket.StatusInternalError, "connection closed")

	for {
		var tickerMsg TickerMessage
		if err := wsjson.Read(context.Background(), conn, &tickerMsg); err != nil {
			log.Println(err)
		}

		ticker := tickerMsg.Data
		err = redisDb.SaveTicker(&ticker, ticker.Sym)
		if err != nil {
			log.Println(err)
		}

		fmt.Print(".")
	}
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
