// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: userfollow.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getFeedId = `-- name: GetFeedId :one
SELECT id
FROM feeds
WHERE url = $1
`

func (q *Queries) GetFeedId(ctx context.Context, url string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getFeedId, url)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
