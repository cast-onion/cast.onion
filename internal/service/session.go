package service

import (
	"context"
	"time"

	"github.com/cast-onion/internal/cache"
	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/pkg/keygen"
)

type SessionService struct {
	q     *queries.Queries
	redis *cache.Redis
}

func NewSessionService(q *queries.Queries, redis *cache.Redis) *SessionService {
	return &SessionService{q: q, redis: redis}
}

const sessionTTL = 24 * time.Hour

func (s *SessionService) Create(ctx context.Context, ip, userAgent string) (string, error) {
	id, err := keygen.GenerateSessionID()
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(sessionTTL)
	if err := s.q.InsertSession(ctx, id, expiresAt, ip, userAgent); err != nil {
		return "", err
	}
	s.redis.SetSession(ctx, id, sessionTTL)
	return id, nil
}

func (s *SessionService) Validate(ctx context.Context, id string) (bool, error) {
	ok, err := s.redis.SessionExists(ctx, id)
	if err == nil && ok {
		return true, nil
	}
	return s.q.SessionValid(ctx, id)
}

func (s *SessionService) Invalidate(ctx context.Context, id string) error {
	s.redis.DeleteSession(ctx, id)
	return s.q.InvalidateSession(ctx, id)
}
