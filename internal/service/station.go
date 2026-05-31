package service

import (
	"context"

	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/internal/model"
	"github.com/cast-onion/pkg/cdn"
	"github.com/google/uuid"
)

type StationService struct {
	q *queries.Queries
}

func NewStationService(q *queries.Queries) *StationService {
	return &StationService{q: q}
}

func resolveArt(st *model.Station) {
	if st.ArtKey != nil && *st.ArtKey != "" {
		url := cdn.AssetURL(*st.ArtKey)
		st.ArtKey = &url
	}
}

func (s *StationService) Directory(ctx context.Context) ([]*model.Station, error) {
	stations, err := s.q.ListActiveStations(ctx)
	if err != nil {
		return nil, err
	}
	for _, st := range stations {
		resolveArt(st)
	}
	return stations, nil
}

func (s *StationService) GetByID(ctx context.Context, id string) (*model.Station, error) {
	st, err := s.q.GetStation(ctx, id)
	if err != nil {
		return nil, err
	}
	resolveArt(st)
	return st, nil
}

func (s *StationService) GetBySlug(ctx context.Context, slug string) (*model.Station, error) {
	st, err := s.q.GetStationBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	resolveArt(st)
	return st, nil
}

func (s *StationService) Create(ctx context.Context, name, description, genre string) (*model.Station, error) {
	id := uuid.NewString()
	base := slugify(name)
	slug := base
	if slug == "" {
		slug = id[:8]
	} else {
		slug = slug + "-" + id[:8]
	}
	st := &model.Station{
		ID:          id,
		Slug:        slug,
		DisplayName: name,
		Description: &description,
		Genre:       &genre,
		Status:      "active",
	}
	if err := s.q.InsertStation(ctx, st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *StationService) Update(ctx context.Context, stationID, description, genre, websiteURL string) error {
	return s.q.UpdateStation(ctx, stationID, description, genre, websiteURL)
}

func (s *StationService) SetStatus(ctx context.Context, stationID, status string) error {
	return s.q.UpdateStationStatus(ctx, stationID, status)
}

func slugify(name string) string {
	slug := make([]byte, 0, len(name))
	for _, c := range []byte(name) {
		switch {
		case c >= 'a' && c <= 'z':
			slug = append(slug, c)
		case c >= 'A' && c <= 'Z':
			slug = append(slug, c+32)
		case c == ' ' || c == '-':
			slug = append(slug, '-')
		}
	}
	return string(slug)
}
