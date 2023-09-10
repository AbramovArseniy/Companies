-- name: GetNodeTags :many
SELECT * FROM tags WHERE node_id = $1;