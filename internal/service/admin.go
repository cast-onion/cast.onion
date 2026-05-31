package service

import (
	"context"
	"errors"

	"github.com/cast-onion/internal/api/websocket"
	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/internal/email"
	"github.com/cast-onion/internal/model"
	"github.com/google/uuid"
)

var ErrAlreadyReviewed = errors.New("application already reviewed")

type Broadcaster interface {
	Broadcast(eventType, stationID string)
}

type AdminService struct {
	q          *queries.Queries
	stationSvc *StationService
	authSvc    *AuthService
	hub        Broadcaster
	mailer     *email.Dispatcher
}

func NewAdminService(q *queries.Queries, stationSvc *StationService, authSvc *AuthService, hub Broadcaster, mailer *email.Dispatcher) *AdminService {
	return &AdminService{q: q, stationSvc: stationSvc, authSvc: authSvc, hub: hub, mailer: mailer}
}

type ApprovalResult struct {
	StationID      string
	RawStationKey  string
	RawAccessToken string
}

func (s *AdminService) ListPending(ctx context.Context) ([]*model.Application, error) {
	return s.q.ListApplicationsByStatus(ctx, "pending")
}

func (s *AdminService) ApproveApplication(ctx context.Context, adminID, appID, reason string) (*ApprovalResult, error) {
	app, err := s.q.GetApplication(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app.Status != "pending" {
		return nil, ErrAlreadyReviewed
	}

	station, err := s.stationSvc.Create(ctx, app.StationName, app.Description, app.Genre)
	if err != nil {
		return nil, err
	}

	creds, err := s.authSvc.IssueCredentials(ctx, station.ID)
	if err != nil {
		return nil, err
	}

	if err := s.q.UpdateApplicationStatus(ctx, appID, "approved", adminID, station.ID); err != nil {
		return nil, err
	}

	s.q.InsertAdminAction(ctx, uuid.NewString(), adminID, "approve", "application", appID, reason)
	s.hub.Broadcast("application_approved", station.ID)

	s.mailer.ApplicationApproved(
		app.ContactEmail,
		app.ContactEmail,
		app.StationName,
		creds.RawStationKey,
		creds.RawAccessToken,
	)

	return &ApprovalResult{
		StationID:      station.ID,
		RawStationKey:  creds.RawStationKey,
		RawAccessToken: creds.RawAccessToken,
	}, nil
}

func (s *AdminService) DenyApplication(ctx context.Context, adminID, appID, reason string) error {
	app, err := s.q.GetApplication(ctx, appID)
	if err != nil {
		return err
	}
	if app.Status != "pending" {
		return ErrAlreadyReviewed
	}

	if err := s.q.UpdateApplicationStatus(ctx, appID, "denied", adminID, ""); err != nil {
		return err
	}

	s.q.InsertAdminAction(ctx, uuid.NewString(), adminID, "deny", "application", appID, reason)
	s.hub.Broadcast("application_denied", "")

	s.mailer.ApplicationDenied(app.ContactEmail, app.ContactEmail, app.StationName, reason)

	return nil
}

func (s *AdminService) SuspendStation(ctx context.Context, adminID, stationID, reason string) error {
	station, err := s.stationSvc.GetByID(ctx, stationID)
	if err != nil {
		return err
	}

	if err := s.stationSvc.SetStatus(ctx, stationID, "suspended"); err != nil {
		return err
	}
	if err := s.authSvc.RevokeAll(ctx, stationID); err != nil {
		return err
	}

	s.q.InsertAdminAction(ctx, uuid.NewString(), adminID, "suspend", "station", stationID, reason)
	s.hub.Broadcast("station_suspended", stationID)

	if email := s.getStationEmail(ctx, stationID); email != "" {
		s.mailer.StationSuspended(email, email, station.DisplayName, reason)
	}

	return nil
}

func (s *AdminService) RevokeStation(ctx context.Context, adminID, stationID, reason string) error {
	station, err := s.stationSvc.GetByID(ctx, stationID)
	if err != nil {
		return err
	}

	if err := s.stationSvc.SetStatus(ctx, stationID, "revoked"); err != nil {
		return err
	}
	if err := s.authSvc.RevokeAll(ctx, stationID); err != nil {
		return err
	}

	s.q.InsertAdminAction(ctx, uuid.NewString(), adminID, "revoke", "station", stationID, reason)
	s.hub.Broadcast("station_revoked", stationID)

	if email := s.getStationEmail(ctx, stationID); email != "" {
		s.mailer.StationRevoked(email, email, station.DisplayName, reason)
	}

	return nil
}

func (s *AdminService) UnsuspendStation(ctx context.Context, adminID, stationID, reason string) error {
	if err := s.stationSvc.SetStatus(ctx, stationID, "active"); err != nil {
		return err
	}
	s.q.InsertAdminAction(ctx, uuid.NewString(), adminID, "unsuspend", "station", stationID, reason)
	s.hub.Broadcast("station_unsuspended", stationID)
	return nil
}

func (s *AdminService) getStationEmail(ctx context.Context, stationID string) string {
	apps, err := s.q.ListAllApplications(ctx)
	if err != nil {
		return ""
	}
	for _, a := range apps {
		if a.StationID != nil && *a.StationID == stationID {
			return a.ContactEmail
		}
	}
	return ""
}

// satisfy websocket import
var _ = websocket.Event{}
