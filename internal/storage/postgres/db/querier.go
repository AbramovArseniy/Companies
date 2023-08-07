// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	GetAllTree(ctx context.Context) ([]GetAllTreeRow, error)
	GetHierarchy(ctx context.Context, id sql.NullInt32) ([]GetHierarchyRow, error)
	GetOneNode(ctx context.Context, id sql.NullInt32) (GetOneNodeRow, error)
}

var _ Querier = (*Queries)(nil)
