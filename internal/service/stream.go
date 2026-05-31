package service

import (
	"context"
	"errors"

	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/pkg/keygen"
)

var ErrUnauthorizedStream = errors.New("stream key invalid or station not active")

type StreamService struct {
	q *queries.Queries
}

func NewStreamService(q *queries.Queries) *StreamService {
	return &StreamService{q: q}
}

func (s *StreamService) ValidateStationKey(ctx context.Context, rawKey string) error {
	hash := keygen.HashRaw(rawKey)
	stationID, err := s.q.GetStationIDByKey(ctx, hash)
	if err != nil {
		return ErrUnauthorizedStream
	}
	station, err := s.q.GetStation(ctx, stationID)
	if err != nil || station.Status != "active" {
		return ErrUnauthorizedStream
	}
	return nil
}
