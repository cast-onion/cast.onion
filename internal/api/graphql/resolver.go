package graphql

import (
	"github.com/cast-onion/internal/api/graphql/resolvers"
	"github.com/cast-onion/internal/service"
	"github.com/graph-gophers/graphql-go"
)

func NewSchema(stationSvc *service.StationService, appSvc *service.ApplicationService, schemaStr string) *graphql.Schema {
	r := &resolvers.Resolver{
		StationSvc: stationSvc,
		AppSvc:     appSvc,
	}
	return graphql.MustParseSchema(schemaStr, r)
}
