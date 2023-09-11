// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	GetAllTree(ctx context.Context) ([]GetAllTreeRow, error)
	GetChangesNum(ctx context.Context) ([]GetChangesNumRow, error)
	GetHierarchy(ctx context.Context, id int32) ([]GetHierarchyRow, error)
	GetNodeTags(ctx context.Context, nodeID pgtype.Int4) ([]GetNodeTagsRow, error)
	GetOneNode(ctx context.Context, id int32) (GetOneNodeRow, error)
	SaveAlert(ctx context.Context, arg SaveAlertParams) error
	SaveChange(ctx context.Context, arg SaveChangeParams) error
	UpdateTag(ctx context.Context, arg UpdateTagParams) error
}

var _ Querier = (*Queries)(nil)
