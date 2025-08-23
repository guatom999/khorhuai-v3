CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS notifications (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       UUID,
  channel       TEXT NOT NULL,       
  recipient     TEXT NOT NULL,        
  template_name TEXT,
  data          JSONB,
  status        TEXT NOT NULL DEFAULT 'queued',  -- queued|sent|failed|cancelled
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notifications_user   ON notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications (status);

CREATE TABLE IF NOT EXISTS delivery_attempts (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  notification_id UUID NOT NULL,
  status          TEXT NOT NULL,      -- sent|failed
  error_message   TEXT,
  provider_raw    JSONB,
  created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE delivery_attempts
  DROP CONSTRAINT IF EXISTS fk_delivery_attempts_notification,
  ADD  CONSTRAINT fk_delivery_attempts_notification
  FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE;
