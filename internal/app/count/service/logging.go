package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingMiddleware struct {
	logger log.Logger
	next   CountService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next CountService) CountService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (lm loggingMiddleware) Count(ctx context.Context) (res int64, err error) {
	defer func() {
		lm.logger.Log("method", "Count", "err", err)
	}()

	return lm.next.Count(ctx)
}
