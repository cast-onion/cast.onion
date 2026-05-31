package resolvers

import (
	"context"

	"github.com/cast-onion/internal/model"
	"github.com/cast-onion/internal/service"
	"github.com/graph-gophers/graphql-go"
)

type Resolver struct {
	StationSvc *service.StationService
	AppSvc     *service.ApplicationService
}

type stationResolver struct {
	s *model.Station
}

func (r *stationResolver) ID() graphql.ID       { return graphql.ID(r.s.ID) }
func (r *stationResolver) Slug() string         { return r.s.Slug }
func (r *stationResolver) DisplayName() string  { return r.s.DisplayName }
func (r *stationResolver) Description() *string { return r.s.Description }
func (r *stationResolver) Genre() *string       { return r.s.Genre }
func (r *stationResolver) WebsiteUrl() *string  { return r.s.WebsiteURL }
func (r *stationResolver) ArtUrl() *string      { return r.s.ArtKey }
func (r *stationResolver) Status() string       { return r.s.Status }
func (r *stationResolver) CreatedAt() string    { return r.s.CreatedAt.String() }

func (r *Resolver) Directory(ctx context.Context) ([]*stationResolver, error) {
	stations, err := r.StationSvc.Directory(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*stationResolver, len(stations))
	for i, s := range stations {
		out[i] = &stationResolver{s}
	}
	return out, nil
}

func (r *Resolver) Station(ctx context.Context, args struct{ ID graphql.ID }) (*stationResolver, error) {
	s, err := r.StationSvc.GetByID(ctx, string(args.ID))
	if err != nil {
		return nil, err
	}
	return &stationResolver{s}, nil
}

func (r *Resolver) StationBySlug(ctx context.Context, args struct{ Slug string }) (*stationResolver, error) {
	s, err := r.StationSvc.GetBySlug(ctx, args.Slug)
	if err != nil {
		return nil, err
	}
	return &stationResolver{s}, nil
}
