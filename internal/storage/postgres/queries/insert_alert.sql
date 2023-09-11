-- name: SaveAlert :exec
INSERT INTO alerts (type, uuid, alert_time, severity, state)
VALUES ($1, $2, to_timestamp($3/1000), $4, $5);