-- +goose Up
ALTER TABLE products ADD COLUMN user_id uuid REFERENCES users(id);

-- +goose Down
ALTER TABLE products DROP COLUMN IF EXISTS user_id;
