-- name: GetProducts :many
SELECT * FROM products p;

-- name: GetProductsImages :many
SELECT * FROM product_images; 

-- name: GetProductsCategories :many
SELECT c.id, pc.product_id, c.name FROM products_category pc
JOIN categories c ON pc.category_id = c.id; 

-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductCategoriesById :many
SELECT c.id, pc.product_id, c.name FROM products_category pc
JOIN categories c ON pc.category_id = c.id 
WHERE product_id = $1;

-- name: GetProductImagesById :many
SELECT * FROM product_images WHERE product_id = $1;

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
