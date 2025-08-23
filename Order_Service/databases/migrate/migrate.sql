CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS orders (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id          UUID,                      
  currency         CHAR(3) NOT NULL DEFAULT 'USD',
  total_price     BIGINT NOT NULL,
  status           TEXT NOT NULL DEFAULT 'pending', 
  shipping_address JSONB,
  billing_address  JSONB,
  created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_user   ON orders (user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);


CREATE TABLE IF NOT EXISTS order_items (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id         UUID NOT NULL,
  product_id       UUID NOT NULL,             
  title            TEXT NOT NULL,
  sku              TEXT,
  quantity         INTEGER NOT NULL,
  unit_price        BIGINT NOT NULL,
  total_price      BIGINT NOT NULL,
  created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE order_items
  DROP CONSTRAINT IF EXISTS fk_order_items_order,
  ADD CONSTRAINT fk_order_items_order
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
