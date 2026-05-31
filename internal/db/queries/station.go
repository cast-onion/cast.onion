package queries

import (
	"context"

	"github.com/cast-onion/internal/model"
)

func (q *Queries) InsertStation(ctx context.Context, s *model.Station) error {
	_, err := q.db.ExecContext(ctx,
		`INSERT INTO stations (id, slug, display_name, description, genre, status) VALUES (?, ?, ?, ?, ?, ?)`,
		s.ID, s.Slug, s.DisplayName, s.Description, s.Genre, s.Status)
	return err
}

func (q *Queries) GetStation(ctx context.Context, id string) (*model.Station, error) {
	var s model.Station
	err := q.db.GetContext(ctx, &s, `SELECT * FROM stations WHERE id = ?`, id)
	return &s, err
}

func (q *Queries) GetStationBySlug(ctx context.Context, slug string) (*model.Station, error) {
	var s model.Station
	err := q.db.GetContext(ctx, &s, `SELECT * FROM stations WHERE slug = ?`, slug)
	return &s, err
}

func (q *Queries) ListActiveStations(ctx context.Context) ([]*model.Station, error) {
	var stations []*model.Station
	err := q.db.SelectContext(ctx, &stations, `SELECT * FROM stations WHERE status = 'active' ORDER BY display_name`)
	return stations, err
}

func (q *Queries) ListAllStations(ctx context.Context) ([]*model.Station, error) {
	var stations []*model.Station
	err := q.db.SelectContext(ctx, &stations, `SELECT * FROM stations ORDER BY created_at DESC`)
	return stations, err
}

func (q *Queries) UpdateStationStatus(ctx context.Context, id, status string) error {
	_, err := q.db.ExecContext(ctx, `UPDATE stations SET status = ? WHERE id = ?`, status, id)
	return err
}

func (q *Queries) UpdateStation(ctx context.Context, id, description, genre, websiteURL string) error {
	_, err := q.db.ExecContext(ctx,
		`UPDATE stations SET description = ?, genre = ?, website_url = ? WHERE id = ?`,
		description, genre, websiteURL, id)
	return err
}
