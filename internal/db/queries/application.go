package queries

import (
	"context"

	"github.com/cast-onion/internal/model"
)

func (q *Queries) InsertApplication(ctx context.Context, a *model.Application) error {
	_, err := q.db.ExecContext(ctx,
		`INSERT INTO applications (id, session_id, contact_email, station_name, description, genre, notes, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.SessionID, a.ContactEmail, a.StationName, a.Description, a.Genre, a.Notes, a.Status)
	return err
}

func (q *Queries) GetApplication(ctx context.Context, id string) (*model.Application, error) {
	var a model.Application
	err := q.db.GetContext(ctx, &a, `SELECT * FROM applications WHERE id = ?`, id)
	return &a, err
}

func (q *Queries) ListApplicationsByStatus(ctx context.Context, status string) ([]*model.Application, error) {
	var apps []*model.Application
	err := q.db.SelectContext(ctx, &apps, `SELECT * FROM applications WHERE status = ? ORDER BY created_at`, status)
	return apps, err
}

func (q *Queries) UpdateApplicationStatus(ctx context.Context, id, status, reviewedBy, stationID string) error {
	_, err := q.db.ExecContext(ctx,
		`UPDATE applications SET status = ?, reviewed_by = ?, reviewed_at = NOW(), station_id = ? WHERE id = ?`,
		status, reviewedBy, nullableStr(stationID), id)
	return err
}

func (q *Queries) ListAllApplications(ctx context.Context) ([]*model.Application, error) {
	var apps []*model.Application
	err := q.db.SelectContext(ctx, &apps, `SELECT * FROM applications ORDER BY created_at DESC`)
	return apps, err
}
