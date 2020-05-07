package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(FooService) FooService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type FooService interface {
	// [method=post,expose=true]
	Foo(ctx context.Context, s string) (res string, err error)
}

// the concrete implementation of service interface
type stubFooService struct {
	logger log.Logger `json:"logger"`
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(logger log.Logger) (s FooService) {
	var svc FooService
	{
		svc = &stubFooService{logger: logger}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Foo
func (fo *stubFooService) Foo(ctx context.Context, s string) (res string, err error) {
	return fmt.Sprintf("%s bar", s), err
}
