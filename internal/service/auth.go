package service

import (
	"context"

	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/pkg/keygen"
	"github.com/google/uuid"
)

type AuthService struct {
	q *queries.Queries
}

func NewAuthService(q *queries.Queries) *AuthService {
	return &AuthService{q: q}
}

type IssuedCredentials struct {
	RawStationKey  string
	RawAccessToken string
}

func (s *AuthService) IssueCredentials(ctx context.Context, stationID string) (*IssuedCredentials, error) {
	rawKey, keyHash, err := keygen.GenerateStationKey()
	if err != nil {
		return nil, err
	}
	rawToken, tokenHash, err := keygen.GenerateAccessToken()
	if err != nil {
		return nil, err
	}

	if err := s.q.InsertStationKey(ctx, uuid.NewString(), stationID, keyHash); err != nil {
		return nil, err
	}
	if err := s.q.InsertAccessToken(ctx, uuid.NewString(), stationID, tokenHash); err != nil {
		return nil, err
	}

	return &IssuedCredentials{RawStationKey: rawKey, RawAccessToken: rawToken}, nil
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, raw string) (string, error) {
	hash := keygen.HashRaw(raw)
	return s.q.GetStationIDByToken(ctx, hash)
}

func (s *AuthService) ValidateStationKey(ctx context.Context, raw string) (string, error) {
	hash := keygen.HashRaw(raw)
	return s.q.GetStationIDByKey(ctx, hash)
}

func (s *AuthService) RevokeAll(ctx context.Context, stationID string) error {
	if err := s.q.RevokeAllKeys(ctx, stationID); err != nil {
		return err
	}
	return s.q.RevokeAllTokens(ctx, stationID)
}
