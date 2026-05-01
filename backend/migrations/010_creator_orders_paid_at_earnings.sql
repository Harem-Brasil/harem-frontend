-- Momento em que o pagamento foi confirmado (agregações de ganhos por período).
-- Índice para sumários por creator_id + paid_at (monetização / agregações de ganhos).

ALTER TABLE creator_orders ADD COLUMN IF NOT EXISTS paid_at TIMESTAMPTZ;

UPDATE creator_orders
SET paid_at = updated_at
WHERE paid_at IS NULL
  AND status IN ('paid', 'fulfilled', 'refunded');

CREATE INDEX IF NOT EXISTS idx_creator_orders_earnings_by_paid_at
    ON creator_orders (creator_id, paid_at)
    WHERE status IN ('paid', 'fulfilled')
      AND paid_at IS NOT NULL;
