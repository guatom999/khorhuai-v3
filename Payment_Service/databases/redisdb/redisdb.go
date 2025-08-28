package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/guatom999/ecommerce-payment-api/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	Rdb *redis.Client
}

type Record struct {
	Status     string         `json:"status"`
	PaymentID  string         `json:"payment_id,omitempty"`
	HTTPStatus int            `json:"http_status,omitempty"`
	Response   map[string]any `json:"response,omitempty"`
}

func NewRedis(cfg *config.Config) *Store {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		log.Printf("Error Instrument Metrics for redis failed %v", err)
	}

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Printf("Error Instrument Metrics for redis failed %v", err)
	}

	return &Store{Rdb: rdb}
}

func (s *Store) TryStart(ctx context.Context, key string, ttl time.Duration) (started bool, rec *Record, err error) {

	b, _ := json.Marshal(&Record{Status: "processing"})
	ok, err := s.Rdb.SetNX(ctx, fmt.Sprintf("idem:%s", key), b, ttl).Result()
	if err != nil {
		return false, nil, err
	}

	if ok {
		return true, &Record{Status: "processing"}, nil
	}

	val, err := s.Rdb.Get(ctx, fmt.Sprintf("idem:%s", key)).Bytes()
	if err == redis.Nil {
		return false, nil, nil
	}

	if err != nil {
		return false, nil, err
	}

	r := new(Record)
	_ = json.Unmarshal(val, r)

	return false, r, nil

}

func (s *Store) Complete(ctx context.Context, key string, paymentID string, httpStatus int, resp map[string]any, ttl time.Duration) error {

	rec := &Record{Status: "done", PaymentID: paymentID, HTTPStatus: httpStatus, Response: resp}
	b, _ := json.Marshal(rec)
	_, err := s.Rdb.Set(ctx, fmt.Sprintf("idem:%s", key), b, ttl).Result()
	if err != nil {
		return err
	}

	return nil

}
