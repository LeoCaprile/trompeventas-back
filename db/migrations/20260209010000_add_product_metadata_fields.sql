-- +goose Up
ALTER TABLE products ADD COLUMN condition TEXT NOT NULL DEFAULT 'Nuevo';
ALTER TABLE products ADD COLUMN state TEXT NOT NULL DEFAULT 'Disponible';
ALTER TABLE products ADD COLUMN negotiable TEXT NOT NULL DEFAULT 'No conversable';

-- +goose Down
ALTER TABLE products DROP COLUMN IF EXISTS condition;
ALTER TABLE products DROP COLUMN IF EXISTS state;
ALTER TABLE products DROP COLUMN IF EXISTS negotiable;
