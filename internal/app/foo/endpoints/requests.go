package endpoints

import (
	"github.com/cage1016/gae-lab-001/internal/app/foo/service"
	"github.com/cage1016/gae-lab-001/internal/pkg/errors"
)

const (
	maxLimitSize = 100
)

type Request interface {
	validate() error
}

// InsertRequest collects the request parameters for the Insert method.
type InsertRequest struct {
	S string `json:"s"`
}

func (r InsertRequest) validate() error {
	return nil // TBA
}

// ListRequest collects the request parameters for the List method.
type ListRequest struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

func (r ListRequest) validate() error {
	if r.Limit <= 0 || r.Limit > maxLimitSize {
		return errors.Wrap(service.ErrMalformedEntity, errors.New("limit must between 1 - 100"))
	}

	return nil
}
