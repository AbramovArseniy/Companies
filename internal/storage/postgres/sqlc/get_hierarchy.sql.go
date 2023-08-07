// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_hierarchy.sql

package db

import (
	"context"
	"database/sql"
)

const getIerarchy = `-- name: GetIerarchy :many
WITH RECURSIVE r(id, name, parent_id, level) AS
        (SELECT tr.id, tr.name, tr.parent_id, 1
        FROM nodes tl
        LEFT JOIN nodes tr 
        ON tl.parent_id = tr.id
        WHERE tl.id = $1
        UNION ALL
        SELECT t.id, t.name, t.parent_id, level+1
        FROM nodes t, r
        WHERE t.id = r.parent_id )
        SELECT id, name, parent_id, ROW_NUMBER() OVER (ORDER BY level DESC) AS level, info.address, info.phone_number, info.contant_person FROM r AS ierarchy 
        LEFT JOIN info 
        ON hierarchy.id = info.node_id
`

type GetIerarchyRow struct {
	ID            sql.NullInt32  `json:"id"`
	Name          sql.NullString `json:"name"`
	ParentID      sql.NullInt32  `json:"parent_id"`
	Level         int64          `json:"level"`
	Address       sql.NullString `json:"address"`
	PhoneNumber   sql.NullString `json:"phone_number"`
	ContantPerson sql.NullString `json:"contant_person"`
}

func (q *Queries) GetIerarchy(ctx context.Context, id sql.NullInt32) ([]GetIerarchyRow, error) {
	rows, err := q.db.QueryContext(ctx, getIerarchy, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetIerarchyRow{}
	for rows.Next() {
		var i GetIerarchyRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.Level,
			&i.Address,
			&i.PhoneNumber,
			&i.ContantPerson,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
