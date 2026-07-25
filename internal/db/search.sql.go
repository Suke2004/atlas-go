package db

import (
	"context"
)

const search = `-- name: Search :many
SELECT entity_type, entity_id, user_id, title, content, tags, rank
FROM search_index
WHERE user_id = ? AND search_index MATCH ?
ORDER BY rank
LIMIT 20
`

type SearchParams struct {
	UserID int64  `json:"user_id"`
	Query  string `json:"query"`
}

func (q *Queries) Search(ctx context.Context, arg SearchParams) ([]SearchIndexRow, error) {
	rows, err := q.db.QueryContext(ctx, search, arg.UserID, arg.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchIndexRow
	for rows.Next() {
		var i SearchIndexRow
		if err := rows.Scan(
			&i.EntityType,
			&i.EntityID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.Tags,
			&i.Rank,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
