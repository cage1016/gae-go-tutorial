package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/cage1016/gae-lab-001/internal/app/foo/model"
	"github.com/cage1016/gae-lab-001/internal/pkg/errors"
)

var (
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrInvalidQueryParams indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrInvalidQueryParams = errors.New("invalid query params")
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(FooService) FooService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type FooService interface {
	// [method=post,expose=true,router=api/foo/insert]
	Insert(ctx context.Context, s string) (res string, err error)
	// [method=get,expose=true,router=api/foo/list]
	List(ctx context.Context, limit uint64, offset uint64) (res model.FooItemPage, err error)
}

// the concrete implementation of service interface
type stubFooService struct {
	logger  log.Logger
	repo    model.FooRepository
	idpNano NanoIdentityProvider
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(repo model.FooRepository, idpNano NanoIdentityProvider, logger log.Logger) (s FooService) {
	var svc FooService
	{
		svc = &stubFooService{repo: repo, idpNano: idpNano, logger: logger}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Insert
func (fo *stubFooService) Insert(ctx context.Context, s string) (res string, err error) {
	value := fmt.Sprintf("%s bar", s)
	nid, _ := fo.idpNano.ID()

	_, err = fo.repo.Insert(ctx, model.Foo{ID: nid, Value: value})
	if err != nil {
		return "", err
	}
	return value, err
}

// Implement the business logic of List
func (fo *stubFooService) List(ctx context.Context, limit uint64, offset uint64) (res model.FooItemPage, err error) {
	return fo.repo.RetrieveAll(ctx, offset, limit)
}
