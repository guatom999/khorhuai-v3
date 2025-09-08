
CREATE EXTENSION
IF NOT EXISTS pgcrypto;

CREATE TABLE
IF NOT EXISTS stock_levels
(
  product_id UUID PRIMARY KEY,
  stock_qty  INTEGER NOT NULL DEFAULT 0,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE
IF NOT EXISTS stock_reservations
(
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid
(),
  order_id    UUID,                 
  user_id     UUID,                
  status      TEXT NOT NULL DEFAULT 'held', 
  expires_at  TIMESTAMP NOT NULL,   
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE
IF NOT EXISTS stock_reservation_items
(
  reservation_id UUID NOT NULL,
  product_id     UUID NOT NULL,
  quantity       INTEGER NOT NULL CHECK
(quantity > 0),
  PRIMARY KEY
(reservation_id, product_id),
  FOREIGN KEY
(reservation_id) REFERENCES stock_reservations
(id) ON
DELETE CASCADE
);

CREATE UNIQUE INDEX
IF NOT EXISTS ux_stock_reservations_order_active
ON stock_reservations
(order_id)
WHERE status IN
('held','committed');

CREATE INDEX
IF NOT EXISTS idx_stock_reservations_status_exp
ON stock_reservations
(status, expires_at);
