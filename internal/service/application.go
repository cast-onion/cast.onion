package service

import (
	"context"

	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/internal/model"
	"github.com/google/uuid"
)

type ApplicationService struct {
	q *queries.Queries
}

func NewApplicationService(q *queries.Queries) *ApplicationService {
	return &ApplicationService{q: q}
}

func (s *ApplicationService) Submit(ctx context.Context, sessionID, email, stationName, description, genre, notes string) (*model.Application, error) {
	app := &model.Application{
		ID:           uuid.NewString(),
		SessionID:    sessionID,
		ContactEmail: email,
		StationName:  stationName,
		Description:  description,
		Genre:        genre,
		Notes:        notes,
		Status:       "pending",
	}
	if err := s.q.InsertApplication(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *ApplicationService) GetByID(ctx context.Context, id string) (*model.Application, error) {
	return s.q.GetApplication(ctx, id)
}

func (s *ApplicationService) ListPending(ctx context.Context) ([]*model.Application, error) {
	return s.q.ListApplicationsByStatus(ctx, "pending")
}
