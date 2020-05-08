package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/cage1016/gae-lab-001/internal/app/count/service"
)

// Endpoints collects all of the endpoints that compose the count service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	CountEndpoint endpoint.Endpoint
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.CountService, logger log.Logger) (ep Endpoints) {
	var countEndpoint endpoint.Endpoint
	{
		method := "count"
		countEndpoint = MakeCountEndpoint(svc)
		countEndpoint = LoggingMiddleware(log.With(logger, "method", method))(countEndpoint)
		ep.CountEndpoint = countEndpoint
	}

	return ep
}

// MakeCountEndpoint returns an endpoint that invokes Count on the service.
// Primarily useful in a server.
func MakeCountEndpoint(svc service.CountService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		res, err := svc.Count(ctx)
		return CountResponse{Res: res}, err
	}
}

// Count implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Count(ctx context.Context) (res int64, err error) {
	resp, err := e.CountEndpoint(ctx, CountRequest{})
	if err != nil {
		return
	}
	response := resp.(CountResponse)
	return response.Res, nil
}
