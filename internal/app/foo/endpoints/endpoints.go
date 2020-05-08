package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/cage1016/gae-lab-001/internal/app/foo/model"
	"github.com/cage1016/gae-lab-001/internal/app/foo/service"
)

// Endpoints collects all of the endpoints that compose the foo service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	InsertEndpoint endpoint.Endpoint `json:""`
	ListEndpoint   endpoint.Endpoint `json:""`
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.FooService, logger log.Logger) (ep Endpoints) {
	var insertEndpoint endpoint.Endpoint
	{
		method := "insert"
		insertEndpoint = MakeInsertEndpoint(svc)
		insertEndpoint = LoggingMiddleware(log.With(logger, "method", method))(insertEndpoint)
		ep.InsertEndpoint = insertEndpoint
	}

	var listEndpoint endpoint.Endpoint
	{
		method := "list"
		listEndpoint = MakeListEndpoint(svc)
		listEndpoint = LoggingMiddleware(log.With(logger, "method", method))(listEndpoint)
		ep.ListEndpoint = listEndpoint
	}

	return ep
}

// MakeInsertEndpoint returns an endpoint that invokes Insert on the service.
// Primarily useful in a server.
func MakeInsertEndpoint(svc service.FooService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(InsertRequest)
		if err := req.validate(); err != nil {
			return InsertResponse{}, err
		}
		res, err := svc.Insert(ctx, req.S)
		return InsertResponse{Res: res}, err
	}
}

// Insert implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Insert(ctx context.Context, s string) (res string, err error) {
	resp, err := e.InsertEndpoint(ctx, InsertRequest{S: s})
	if err != nil {
		return
	}
	response := resp.(InsertResponse)
	return response.Res, nil
}

// MakeListEndpoint returns an endpoint that invokes List on the service.
// Primarily useful in a server.
func MakeListEndpoint(svc service.FooService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		if err := req.validate(); err != nil {
			return ListResponse{}, err
		}
		res, err := svc.List(ctx, req.Limit, req.Offset)
		return ListResponse{Res: res}, err
	}
}

// List implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) List(ctx context.Context, limit uint64, offset uint64) (res model.FooItemPage, err error) {
	resp, err := e.ListEndpoint(ctx, ListRequest{Limit: limit, Offset: offset})
	if err != nil {
		return
	}
	response := resp.(ListResponse)
	return response.Res, nil
}
