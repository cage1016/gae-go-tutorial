package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/gae-lab-001/internal/app/count/service"
	"github.com/cage1016/gae-lab-001/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*CountResponse)(nil)

	_ httptransport.StatusCoder = (*CountResponse)(nil)
)

// CountResponse collects the response values for the Count method.
type CountResponse struct {
	Res int64 `json:"res"`
	Err error `json:"-"`
}

func (r CountResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r CountResponse) Headers() http.Header {
	return http.Header{}
}

func (r CountResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r}
}
