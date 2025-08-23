CREATE EXTENSION
IF NOT EXISTS pgcrypto;

CREATE TABLE
IF NOT EXISTS payments
(
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid
(),
  order_id         UUID NOT NULL,           
  user_id          UUID,                    
  amount           BIGINT NOT NULL,
  currency         CHAR
(3) NOT NULL DEFAULT 'USD',
  status           TEXT NOT NULL DEFAULT 'processing',
  idempotency_key  TEXT UNIQUE,               
  created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX
IF NOT EXISTS idx_payments_order  ON payments
(order_id);
CREATE INDEX
IF NOT EXISTS idx_payments_user   ON payments
(user_id);
CREATE INDEX
IF NOT EXISTS idx_payments_status ON payments
(status);
