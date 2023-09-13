-- name: GetNodeTags :many
SELECT * FROM tags LEFT JOIN alerts ON tags.uuid=alerts.uuid WHERE node_id = $1 AND (alerts.alert_time = (SELECT MAX(alert_time) FROM alerts WHERE uuid=tags.uuid) OR alerts.uuid IS NULL);