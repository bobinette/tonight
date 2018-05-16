package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

type TagReader struct {
	db *sql.DB
}

func NewTagReader(db *sql.DB) *TagReader {
	return &TagReader{db: db}
}

func (r *TagReader) Tags(ctx context.Context, user, q string) ([]string, error) {
	query := `
SELECT tags.tag, COUNT(*) AS c
FROM tags
JOIN tasks ON tags.task_id = tasks.id
JOIN user_has_tasks ON tasks.id = user_has_tasks.task_id
JOIN users ON user_has_tasks.user_id = users.id
WHERE users.username = ?%s
GROUP BY tags.tag
ORDER BY c DESC, tags.tag
`
	args := []interface{}{user}

	tagCondition := ""
	if q != "" {
		tagCondition = " AND tags.tag LIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", q))
	}
	query = fmt.Sprintf(query, tagCondition)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		var c uint
		if err := rows.Scan(&tag, &c); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return tags, nil
}
