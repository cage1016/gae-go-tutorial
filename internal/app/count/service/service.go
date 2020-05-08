package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/gomodule/redigo/redis"

	"github.com/cage1016/gae-lab-001/internal/pkg/errors"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(CountService) CountService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type CountService interface {
	// [method=get,expose=true]
	Count(ctx context.Context) (res int64, err error)
}

// the concrete implementation of service interface
type stubCountService struct {
	logger log.Logger
	conn   redis.Conn
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(conn redis.Conn, logger log.Logger) (s CountService) {
	var svc CountService
	{
		svc = &stubCountService{conn: conn, logger: logger}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Count
func (co *stubCountService) Count(ctx context.Context) (res int64, err error) {
	counter, err := redis.Int(co.conn.Do("INCR", "visits"))
	if err != nil {
		return 0, errors.New("Error incrementing visitor counter")
	}

	return int64(counter), nil
}
