// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_changes_number.sql

package db

import (
	"context"
)

const getChangesNum = `-- name: GetChangesNum :many
SELECT tags.name, COUNT(*) FROM tags_journal  RIGHT JOIN tags ON tags_journal.uuid = tags.uuid WHERE NOW() - interval '5 minute' <= tags_journal.change_time GROUP BY tags.name
`

type GetChangesNumRow struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (q *Queries) GetChangesNum(ctx context.Context) ([]GetChangesNumRow, error) {
	rows, err := q.db.Query(ctx, getChangesNum)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetChangesNumRow
	for rows.Next() {
		var i GetChangesNumRow
		if err := rows.Scan(&i.Name, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
