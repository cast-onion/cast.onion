package queries

import (
	"context"
	"time"
)

func (q *Queries) InsertSession(ctx context.Context, id string, expiresAt time.Time, ip, userAgent string) error {
	_, err := q.db.ExecContext(ctx,
		`INSERT INTO sessions (id, expires_at, ip_address, user_agent) VALUES (?, ?, ?, ?)`,
		id, expiresAt, ip, userAgent)
	return err
}

func (q *Queries) SessionValid(ctx context.Context, id string) (bool, error) {
	var count int
	err := q.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM sessions WHERE id = ? AND expires_at > NOW() AND invalidated = 0`, id).Scan(&count)
	return count > 0, err
}

func (q *Queries) InvalidateSession(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, `UPDATE sessions SET invalidated = 1 WHERE id = ?`, id)
	return err
}
