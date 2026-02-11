-- name: GetProducts :many
SELECT * FROM products p;

-- name: SearchProducts :many
SELECT * FROM products p
WHERE p.name ILIKE '%' || $1 || '%'
   OR p.description ILIKE '%' || $1 || '%';

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
(name, description, price, user_id, condition, state, negotiable)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteProduct :many
DELETE FROM products WHERE id = $1
RETURNING *;

-- name: DeleteProductByOwner :exec
DELETE FROM products WHERE id = $1 AND user_id = $2;

-- name: UpdateProduct :many
UPDATE products
SET name=coalesce(sqlc.narg('name'), name), description=coalesce(sqlc.narg('description'),description), price=coalesce(sqlc.narg('price'), price), condition=coalesce(sqlc.narg('condition'), condition), state=coalesce(sqlc.narg('state'), state), negotiable=coalesce(sqlc.narg('negotiable'), negotiable), updated_at=NOW()
WHERE id=$1
RETURNING *;

-- name: GetProductsByUserId :many
SELECT * FROM products WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateProductImage :one
INSERT INTO product_images (product_id, image_url) VALUES ($1, $2) RETURNING *;

-- name: CreateProductCategory :one
INSERT INTO products_category (product_id, category_id) VALUES ($1, $2) RETURNING *;
