CREATE EXTENSION
IF NOT EXISTS pgcrypto;

CREATE TABLE
IF NOT EXISTS carts
(
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid
(),
  user_id    UUID,
  session_id TEXT,
  currency   CHAR
(3) NOT NULL DEFAULT 'BTH',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT carts_owner_xor CHECK
((user_id IS NULL) <>
(session_id IS NULL))
);

CREATE UNIQUE INDEX
IF NOT EXISTS ux_carts_user
  ON carts
(user_id) WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX
IF NOT EXISTS ux_carts_session
  ON carts
(session_id) WHERE session_id IS NOT NULL;

CREATE TABLE
IF NOT EXISTS cart_items
(
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid
(),
  cart_id          UUID NOT NULL,
  product_id       UUID NOT NULL,
  quantity         INTEGER NOT NULL CHECK
(quantity > 0),
  unit_price_cents BIGINT NOT NULL,
  currency         CHAR
(3) NOT NULL,
  created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE cart_items
  DROP CONSTRAINT IF EXISTS fk_cart_items_cart
,
ADD CONSTRAINT fk_cart_items_cart
  FOREIGN KEY
(cart_id) REFERENCES carts
(id) ON
DELETE CASCADE;

CREATE UNIQUE INDEX
IF NOT EXISTS ux_cart_items_cart_product
  ON cart_items
(cart_id, product_id);

CREATE INDEX
IF NOT EXISTS idx_cart_items_cart ON cart_items
(cart_id);
