-- name: GetChangesNum :one
SELECT COUNT(*) FROM tags_journal WHERE uuid=$1 AND change_time >= now() - interval '5 minute';