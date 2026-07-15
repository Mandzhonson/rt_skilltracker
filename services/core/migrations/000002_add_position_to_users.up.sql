ALTER TABLE users
ADD COLUMN position VARCHAR(150) NOT NULL DEFAULT 'Не указана';

CREATE INDEX idx_users_position ON users(position);