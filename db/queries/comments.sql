-- name: GetCommentsByProductId :many
SELECT
    c.id,
    c.product_id,
    c.user_id,
    c.parent_id,
    c.content,
    c.created_at,
    c.updated_at,
    u.name AS author_name,
    u.image AS author_image,
    (SELECT COUNT(*) FROM comment_votes cv WHERE cv.comment_id = c.id AND cv.vote_type = 'like')::bigint AS likes,
    (SELECT COUNT(*) FROM comment_votes cv WHERE cv.comment_id = c.id AND cv.vote_type = 'dislike')::bigint AS dislikes
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.product_id = $1
ORDER BY c.created_at ASC;

-- name: GetCommentsByProductIdWithUserVote :many
SELECT
    c.id,
    c.product_id,
    c.user_id,
    c.parent_id,
    c.content,
    c.created_at,
    c.updated_at,
    u.name AS author_name,
    u.image AS author_image,
    (SELECT COUNT(*) FROM comment_votes cv WHERE cv.comment_id = c.id AND cv.vote_type = 'like')::bigint AS likes,
    (SELECT COUNT(*) FROM comment_votes cv WHERE cv.comment_id = c.id AND cv.vote_type = 'dislike')::bigint AS dislikes,
    COALESCE((SELECT cv.vote_type FROM comment_votes cv WHERE cv.comment_id = c.id AND cv.user_id = $2), '')::text AS user_vote
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.product_id = $1
ORDER BY c.created_at ASC;

-- name: CreateComment :one
INSERT INTO comments (product_id, user_id, parent_id, content)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1 AND user_id = $2;

-- name: UpsertCommentVote :one
INSERT INTO comment_votes (comment_id, user_id, vote_type)
VALUES ($1, $2, $3)
ON CONFLICT (comment_id, user_id)
DO UPDATE SET vote_type = $3
RETURNING *;

-- name: DeleteCommentVote :exec
DELETE FROM comment_votes WHERE comment_id = $1 AND user_id = $2;
