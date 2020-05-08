package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/gae-lab-001/internal/app/foo/model"
	"github.com/cage1016/gae-lab-001/internal/app/foo/service"
	"github.com/cage1016/gae-lab-001/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*InsertResponse)(nil)

	_ httptransport.StatusCoder = (*InsertResponse)(nil)

	_ httptransport.Headerer = (*ListResponse)(nil)

	_ httptransport.StatusCoder = (*ListResponse)(nil)
)

// InsertResponse collects the response values for the Insert method.
type InsertResponse struct {
	Res string `json:"res"`
	Err error  `json:"-"`
}

func (r InsertResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r InsertResponse) Headers() http.Header {
	return http.Header{}
}

func (r InsertResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r}
}

// ListResponse collects the response values for the List method.
type ListResponse struct {
	Res model.FooItemPage `json:"res"`
	Err error             `json:"-"`
}

func (r ListResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r ListResponse) Headers() http.Header {
	return http.Header{}
}

func (r ListResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r.Res}
}
