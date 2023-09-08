-- name: UpdateTag :exec
UPDATE tags
SET value = $2
WHERE uuid = $1;