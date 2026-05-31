package queries

import (
	"context"

	"github.com/cast-onion/internal/model"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) InsertAdminAction(ctx context.Context, id, adminID, action, targetType, targetID, reason string) error {
	_, err := q.db.ExecContext(ctx,
		`INSERT INTO admin_actions (id, admin_id, action, target_type, target_id, reason) VALUES (?, ?, ?, ?, ?, ?)`,
		id, adminID, action, targetType, targetID, reason)
	return err
}

func (q *Queries) GetAdminIDByTokenHash(ctx context.Context, hash string) (string, error) {
	var adminID string
	err := q.db.QueryRowContext(ctx,
		`SELECT id FROM users WHERE password_hash = ? AND role = 'admin'`, hash).Scan(&adminID)
	return adminID, err
}

func nullableStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (q *Queries) ListAdminActions(ctx context.Context) ([]*model.AdminAction, error) {
	var actions []*model.AdminAction
	err := q.db.SelectContext(ctx, &actions, `SELECT * FROM admin_actions ORDER BY created_at DESC`)
	return actions, err
}
