-- name: GetChangesNum :many
SELECT tags.name, COUNT(*) FROM tags_journal  RIGHT JOIN tags ON tags_journal.uuid = tags.uuid GROUP BY tags.name;