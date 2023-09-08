// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
)

type Querier interface {
	GetAllTree(ctx context.Context) ([]GetAllTreeRow, error)
	GetChangesNum(ctx context.Context, uuid string) (int64, error)
	GetHierarchy(ctx context.Context, id int32) ([]GetHierarchyRow, error)
	GetOneNode(ctx context.Context, id int32) (GetOneNodeRow, error)
	SaveChange(ctx context.Context, arg SaveChangeParams) error
	UpdateTag(ctx context.Context, arg UpdateTagParams) error
}

var _ Querier = (*Queries)(nil)
