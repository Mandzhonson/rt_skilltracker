ALTER TABLE plans
ADD COLUMN generation_status VARCHAR(20)
NOT NULL DEFAULT 'pending'
CHECK (generation_status IN ('pending', 'processing', 'ready', 'failed'));