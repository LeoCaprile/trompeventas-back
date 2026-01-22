-- name: GetCategories :many
SELECT * FROM categories;

-- name: CreateCategory :one
INSERT INTO categories
(name)
VALUES($1)
RETURNING *;

-- name: DeleteCategory :many
DELETE FROM categories WHERE id = $1
RETURNING *;

-- name: UpdateCategory :many
UPDATE categories
SET name=coalesce(sqlc.narg('name'), name)
WHERE id=$1
RETURNING *;

