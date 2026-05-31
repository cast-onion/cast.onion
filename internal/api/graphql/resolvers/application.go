package resolvers

import (
	"context"

	"github.com/cast-onion/internal/model"
	"github.com/graph-gophers/graphql-go"
)

type applicationStatusResolver struct {
	a *model.Application
}

func (r *applicationStatusResolver) ID() graphql.ID     { return graphql.ID(r.a.ID) }
func (r *applicationStatusResolver) Status() string     { return r.a.Status }
func (r *applicationStatusResolver) StationId() *string { return r.a.StationID }

func (r *Resolver) ApplicationStatus(ctx context.Context, args struct{ ID graphql.ID }) (*applicationStatusResolver, error) {
	app, err := r.AppSvc.GetByID(ctx, string(args.ID))
	if err != nil {
		return nil, err
	}
	return &applicationStatusResolver{app}, nil
}
