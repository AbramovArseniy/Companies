-- name: SaveChange :exec
INSERT INTO tags_journal (uuid, change_time)
VALUES ($1, to_timestamp($2/1000));