package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/guatom999/ecommerce-payment-api/config"
	"github.com/guatom999/ecommerce-payment-api/databases"
	"github.com/guatom999/ecommerce-payment-api/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type OutboxRow struct {
	ID            string    `db:"id"`
	Topic         string    `db:"topic"`
	Key           string    `db:"key"`
	Payload       []byte    `db:"payload"`
	RetryCount    int       `db:"retry_count"`
	NextAttemptAt time.Time `db:"next_attempt_at"`
	CreatedAt     time.Time `db:"created_at"`
}

func main() {
	cfg := config.NewConfig()

	brokers := cfg.Kafka.Brokers
	batch := cfg.Outbox.Batch
	interval := utils.Duration(cfg.Outbox.Interval)
	maxRetry := cfg.Outbox.MaxRetry

	db := databases.ConnDB(cfg)
	defer db.Close()

	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Balancer: &kafka.Hash{},
	}
	defer w.Close()

	log.Printf("[outbox] start relay brokers=%s batch=%d interval=%s", brokers, batch, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		ctx := context.Background()
		rows, err := fetchPending(ctx, db, batch)
		if err != nil {
			log.Printf("[outbox] fetch error: %v", err)
			continue
		}
		if len(rows) == 0 {
			continue
		}

		for _, r := range rows {
			msg := kafka.Message{
				Topic: r.Topic,
				Key:   []byte(r.Key),
				Value: r.Payload,
				Time:  time.Now(),
				// Headers: []kafka.Header{{Key:"content-type", Value:[]byte("application/json")}} Not used for now,
			}
			if err := w.WriteMessages(ctx, msg); err != nil {
				log.Printf("[outbox] publish fail id=%s err=%v", r.ID, err)
				_ = markFail(ctx, db, r.ID, r.RetryCount, maxRetry)
				continue
			}
			_ = markSent(ctx, db, r.ID)
		}
	}
}

// ---------- DB helpers ----------

func fetchPending(ctx context.Context, db *sqlx.DB, limit int) ([]OutboxRow, error) {
	rows := []OutboxRow{}
	err := db.SelectContext(ctx, &rows, `
	  WITH cte AS (
	    SELECT id, topic, key, payload, retry_count, next_attempt_at, created_at
	    FROM outbox_events
	    WHERE status='pending' AND next_attempt_at <= CURRENT_TIMESTAMP
	    ORDER BY created_at
	    LIMIT $1
	    FOR UPDATE SKIP LOCKED
	  )
	  UPDATE outbox_events o
	  SET status='processing'
	  FROM cte
	  WHERE o.id = cte.id
	  RETURNING o.id, o.topic, o.key, o.payload, o.retry_count, o.next_attempt_at, o.created_at
	`, limit)
	return rows, err
}

func markSent(ctx context.Context, db *sqlx.DB, id string) error {
	_, err := db.ExecContext(ctx, `
	  UPDATE outbox_events
	  SET status='sent', sent_at=CURRENT_TIMESTAMP
	  WHERE id=$1
	`, id)
	return err
}

func markFail(ctx context.Context, db *sqlx.DB, id string, retryCount, maxRetry int) error {
	retry := retryCount + 1
	status := "pending"
	delay := backoff(retry) // seconds
	if retry > maxRetry {
		status = "failed"
		delay = 300 // 5 min cool down
	}
	_, err := db.ExecContext(ctx, `
	  UPDATE outbox_events
	  SET status=$2,
	      retry_count=$3,
	      next_attempt_at = CURRENT_TIMESTAMP +  (secs => $4)
	  WHERE id=$1
	`, id, status, retry, delay)
	return err
}

func backoff(n int) int {
	// sec := 1 << (n - 1)
	sec := int(math.Pow(2, float64(n-1)))
	if sec > 60 {
		sec = 60
	}
	return sec
}
