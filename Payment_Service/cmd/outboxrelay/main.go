package main

import (
	"context"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/guatom999/ecommerce-payment-api/config"
	"github.com/guatom999/ecommerce-payment-api/databases"
	"github.com/guatom999/ecommerce-payment-api/modules/metrics"
	"github.com/guatom999/ecommerce-payment-api/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
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

	ctx := context.Background()

	_ = ctx

	utils.InitLogger()
	shutdown, _ := utils.InitTracing(ctx, utils.OtelConfig{
		ServiceName: "payment-outboxrelay",
		Endpoint:    cfg.Otel.Endpoint,
		SampleRatio: 1.0,
	})
	defer shutdown(ctx)

	metrics.MustRegister()
	prometheus.MustRegister()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":2112", nil)
	}()

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
		rows, err := fetchPending(ctx, db, batch)
		if err != nil {
			log.Printf("[outbox] fetch error: %v", err)
			utils.AppLogger().Error("fetch pending", zap.Error(err))
			continue
		}
		if len(rows) == 0 {
			continue
		}

		for _, r := range rows {

			start := time.Now()
			_ = start
			_, span := otel.Tracer("outbox").Start(ctx, "outbox.publish")
			span.SetAttributes(attribute.String("topic", r.Topic), attribute.String("key", r.Key))

			msg := kafka.Message{
				Topic: r.Topic,
				Key:   []byte(r.Key),
				Value: r.Payload,
				Time:  time.Now(),
				// Headers: []kafka.Header{{Key:"content-type", Value:[]byte("application/json")}} Not used for now,
			}
			if err := w.WriteMessages(ctx, msg); err != nil {
				metrics.OutboxFailed.Inc()
				span.RecordError(err)
				span.End()
				log.Printf("[outbox] publish fail id=%s err=%v", r.ID, err)
				_ = markFail(ctx, db, r.ID, r.RetryCount, maxRetry)
				continue
			}
			span.End()
			metrics.OutboxPublished.Inc()
			metrics.OutboxLatency.Observe(time.Since(start).Seconds())
			_ = markSent(ctx, db, r.ID)
		}
	}
}

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
