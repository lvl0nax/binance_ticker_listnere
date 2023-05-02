package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Storage interface {
	GetPairsMapping() (map[string]string, error)
	GetTicker(string) (*Ticker, error)
	SaveTicker(*Ticker, string) error
	Set(string, string) error
	Get(string) (string, error)
	Close() error
}

type RedisStorage struct {
	redis *redis.Client
}

func NewStorage(connectionUrl string) Storage {
	opt, err := redis.ParseURL(connectionUrl)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	return &RedisStorage{
		redis: rdb,
	}
}

// TODO: check if needed, because it's not easy to work with Rails cache structure
func (s *RedisStorage) GetPairsMapping() (map[string]string, error) {
	//return s.redis.HGetAll(ctx, "pairs").Result()
	return make(map[string]string), nil
}

func (s *RedisStorage) Close() error {
	return s.redis.Close()
}

func (s *RedisStorage) GetTicker(pair string) (*Ticker, error) {
	ticker := &Ticker{}
	err := s.redis.Get(ctx, pair).Scan(ticker)
	if err != nil {
		return nil, err
	}

	return ticker, nil
}

func (s *RedisStorage) SaveTicker(ticker *Ticker, pair string) error {
	//return s.redis.Set(ctx, pair, ticker, 0).Err()
	return s.redis.HSet(
		ctx,
		pair,
		"last",
		fmt.Sprint(ticker.Last),
		"bid",
		ticker.Bid,
		"ask",
		ticker.Ask,
		"volume",
		ticker.Volume,
		"PriceChangePercent",
		ticker.PriceChangePercent,
	).Err()
}

func (s *RedisStorage) Set(key string, value string) error {
	return s.redis.Set(ctx, key, value, time.Minute).Err()
}

func (s *RedisStorage) Get(key string) (string, error) {
	return s.redis.Get(ctx, key).Result()
}
