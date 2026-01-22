-- name: GetProducts :many
SELECT * FROM products;

-- name: CreateProduct :one
INSERT INTO products
(name, description, price)
VALUES($1, $2, $3)
RETURNING *;

-- name: DeleteProduct :many
DELETE FROM products WHERE id = $1
RETURNING *;

-- name: UpdateProduct :many
UPDATE products
SET name=coalesce(sqlc.narg('name'), name), description=coalesce(sqlc.narg('description'),description), price=coalesce(sqlc.narg('price'), price), updated_at=NOW()
WHERE id=$1
RETURNING *;
