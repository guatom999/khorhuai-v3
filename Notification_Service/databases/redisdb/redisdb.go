package redisDb

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// type Store struct {
// 	Rdb *redis.Client
// }

// type Record struct {
// 	Status     string         `json:"status"`
// 	PaymentID  string         `json:"payment_id,omitempty"`
// 	HTTPStatus int            `json:"http_status,omitempty"`
// 	Response   map[string]any `json:"response,omitempty"`
// }

// func (s *Store) TryStart(ctx context.Context, key string, ttl time.Duration) (started bool, rec *Record, err error) {

// 	b, _ := json.Marshal(&Record{Status: "processing"})
// 	ok, err := s.Rdb.SetNX(ctx, fmt.Sprintf("idem:%s", key), b, ttl).Result()
// 	if err != nil {
// 		return false, nil, err
// 	}

// 	if ok {
// 		return true, &Record{Status: "processing"}, nil
// 	}

// 	val, err := s.Rdb.Get(ctx, fmt.Sprintf("idem:%s", key)).Bytes()
// 	if err == redis.Nil {
// 		return false, nil, nil
// 	}

// 	if err != nil {
// 		return false, nil, err
// 	}

// 	r := new(Record)
// 	_ = json.Unmarshal(val, r)

// 	return false, r, nil

// }

// func (s *Store) Complete(ctx context.Context, key string, paymentID string, httpStatus int, resp map[string]any, ttl time.Duration) error {

// 	rec := &Record{Status: "done", PaymentID: paymentID, HTTPStatus: httpStatus, Response: resp}
// 	b, _ := json.Marshal(rec)
// 	_, err := s.Rdb.Set(ctx, fmt.Sprintf("idem:%s", key), b, ttl).Result()
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
