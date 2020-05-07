package endpoints

import (
	"context"

	"github.com/cage1016/gae-lab-001/internal/app/foo/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Endpoints collects all of the endpoints that compose the foo service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	FooEndpoint endpoint.Endpoint `json:""`
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.FooService, logger log.Logger) (ep Endpoints) {
	var fooEndpoint endpoint.Endpoint
	{
		method := "foo"
		fooEndpoint = MakeFooEndpoint(svc)
		fooEndpoint = LoggingMiddleware(log.With(logger, "method", method))(fooEndpoint)
		ep.FooEndpoint = fooEndpoint
	}

	return ep
}

// MakeFooEndpoint returns an endpoint that invokes Foo on the service.
// Primarily useful in a server.
func MakeFooEndpoint(svc service.FooService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FooRequest)
		if err := req.validate(); err != nil {
			return FooResponse{}, err
		}
		res, err := svc.Foo(ctx, req.S)
		return FooResponse{Res: res}, err
	}
}

// Foo implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Foo(ctx context.Context, s string) (res string, err error) {
	resp, err := e.FooEndpoint(ctx, FooRequest{S: s})
	if err != nil {
		return
	}
	response := resp.(FooResponse)
	return response.Res, nil
}
