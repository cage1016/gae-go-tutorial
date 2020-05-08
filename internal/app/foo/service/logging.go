package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/cage1016/gae-lab-001/internal/app/foo/model"
)

type loggingMiddleware struct {
	logger log.Logger
	next   FooService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next FooService) FooService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (lm loggingMiddleware) Insert(ctx context.Context, s string) (res string, err error) {
	defer func() {
		lm.logger.Log("method", "Insert", "s", s, "err", err)
	}()

	return lm.next.Insert(ctx, s)
}

func (lm loggingMiddleware) List(ctx context.Context, limit uint64, offset uint64) (res model.FooItemPage, err error) {
	defer func() {
		lm.logger.Log("method", "List", "limit", limit, "offset", offset, "err", err)
	}()

	return lm.next.List(ctx, limit, offset)
}
