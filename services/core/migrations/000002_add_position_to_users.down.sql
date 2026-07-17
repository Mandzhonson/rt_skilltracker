DROP INDEX IF EXISTS idx_users_position;

ALTER TABLE users
DROP COLUMN position;