// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_tree.sql

package db

import (
	"context"
	"database/sql"
)

const getAllTree = `-- name: GetAllTree :many
SELECT nodes.id, nodes.name, nodes.parent_id, info.address, info.phone_number, info.contact_person FROM nodes LEFT JOIN info ON nodes.id = info.node_id
`

type GetAllTreeRow struct {
	ID            int32          `json:"id"`
	Name          string         `json:"name"`
	ParentID      sql.NullInt32  `json:"parent_id"`
	Address       sql.NullString `json:"address"`
	PhoneNumber   sql.NullString `json:"phone_number"`
	ContactPerson sql.NullString `json:"contact_person"`
}

func (q *Queries) GetAllTree(ctx context.Context) ([]GetAllTreeRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllTree)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllTreeRow{}
	for rows.Next() {
		var i GetAllTreeRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.Address,
			&i.PhoneNumber,
			&i.ContactPerson,
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
