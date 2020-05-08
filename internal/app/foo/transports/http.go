package transports

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"

	"github.com/cage1016/gae-lab-001/internal/app/foo/endpoints"
	"github.com/cage1016/gae-lab-001/internal/app/foo/service"
	"github.com/cage1016/gae-lab-001/internal/pkg/errors"
	"github.com/cage1016/gae-lab-001/internal/pkg/responses"
)

const (
	contentType string = "application/json"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoints.Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(httpEncodeError),
		httptransport.ServerErrorLogger(logger),
	}

	m := bone.New()
	m.Post("/api/foo/insert", httptransport.NewServer(
		endpoints.InsertEndpoint,
		decodeHTTPInsertRequest,
		encodeJSONResponse,
		append(options, httptransport.ServerBefore(kitjwt.HTTPToContext()))...,
	))
	m.Get("/api/foo/list", httptransport.NewServer(
		endpoints.ListEndpoint,
		decodeHTTPListRequest,
		encodeJSONResponse,
		append(options, httptransport.ServerBefore(kitjwt.HTTPToContext()))...,
	))
	return m
}

// decodeHTTPInsertRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPInsertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.InsertRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// decodeHTTPListRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ListRequest
	var err error
	req.Offset, err = readUintQuery(r, "offset", 0)
	if err != nil {
		return nil, err
	}

	req.Limit, err = readUintQuery(r, "limit", 10)
	if err != nil {
		return nil, err
	}

	return req, nil
}
func httpEncodeError(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	var message string
	var errs []errors.Errors
	w.Header().Set("Content-Type", contentType)

	// HTTP
	switch errorVal := err.(type) {
	case errors.Error:
		switch {
		case errors.Contains(errorVal, service.ErrInvalidQueryParams),
			errors.Contains(errorVal, service.ErrMalformedEntity):
			code = http.StatusBadRequest
		}

		if errorVal.Msg() != "" {
			message, errs = errorVal.Msg(), errorVal.Errors()
		}
	default:
		switch err {
		case io.ErrUnexpectedEOF, io.EOF:
			code = http.StatusBadRequest
		case kitjwt.ErrTokenContextMissing:
			code = http.StatusUnauthorized
		default:
			switch err.(type) {
			case *json.SyntaxError, *json.UnmarshalTypeError:
				code = http.StatusBadRequest
			}
		}

		errs = errors.FromError(err.Error())
		message = errs[0].Message
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(responses.ErrorRes{responses.ErrorResItem{code, message, errs}})
}

func encodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if headerer, ok := response.(httptransport.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	if ar, ok := response.(responses.Responser); ok {
		return json.NewEncoder(w).Encode(ar.Response())
	}

	return json.NewEncoder(w).Encode(response)
}

func readUintQuery(r *http.Request, key string, def uint64) (uint64, error) {
	vals := bone.GetQuery(r, key)
	if len(vals) > 1 {
		return 0, service.ErrInvalidQueryParams
	}

	if len(vals) == 0 {
		return def, nil
	}

	strval := vals[0]
	val, err := strconv.ParseUint(strval, 10, 64)
	if err != nil {
		return 0, service.ErrInvalidQueryParams
	}

	return val, nil
}
