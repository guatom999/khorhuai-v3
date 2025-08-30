-- จำเป็นสำหรับ gen_random_uuid()
CREATE EXTENSION
IF NOT EXISTS pgcrypto;

-- ตารางสินค้า (ของคุณมีอยู่แล้ว)
-- products(id UUID PK, name TEXT, price BIGINT, stock_qty INT, created_at, updated_at)

-- โต้ะจองสต็อก (หัว)
CREATE TABLE
IF NOT EXISTS stock_reservations
(
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid
(),
  order_id    UUID,                 -- อ้างคำสั่งซื้อ (ถ้ารู้ตั้งแต่ตอนจอง)
  user_id     UUID,                 -- เผื่อ track ผู้ใช้
  status      TEXT NOT NULL DEFAULT 'held',  -- held | released | committed
  expires_at  TIMESTAMP NOT NULL,   -- TTL เช่น now() + interval '15 minutes'
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- รายการสินค้าที่จองในแต่ละ reservation
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

-- กันการจองซ้ำซ้อนต่อ order (ทางเลือก: เปิดเฉพาะสถานะ active)
CREATE UNIQUE INDEX
IF NOT EXISTS ux_stock_reservations_order_active
ON stock_reservations
(order_id)
WHERE status IN
('held','committed');

-- ใช้บ่อย
CREATE INDEX
IF NOT EXISTS idx_stock_reservations_status_exp
ON stock_reservations
(status, expires_at);
