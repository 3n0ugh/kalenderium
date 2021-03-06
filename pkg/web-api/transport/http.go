package transport

import (
	"context"
	"encoding/json"
	context2 "github.com/3n0ugh/kalenderium/internal/context"
	errs "github.com/3n0ugh/kalenderium/internal/err"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/3n0ugh/kalenderium/pkg/web-api/endpoints"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	r := mux.NewRouter()

	r.MethodNotAllowedHandler = http.HandlerFunc(errs.MethodNotAllowedResponse)
	r.NotFoundHandler = http.HandlerFunc(errs.NotFoundResponse)

	r.Use(authentication, rateLimit, secureHeaders, enableCORS, recoverPanic, prometheusMiddleware)

	r.Handle("/v1/metrics", promhttp.Handler())
	r.Handle("/v1/calendar", requireAuthenticatedUser(httpTransport.NewServer(
		ep.AddEventEndpoint,
		decodeHTTPAddEventRequest,
		encodeResponse))).Methods(http.MethodPost)

	r.Handle("/v1/calendar", requireAuthenticatedUser(httpTransport.NewServer(
		ep.ListEventEndpoint,
		decodeHTTPListEventRequest,
		encodeResponse))).Methods(http.MethodGet)

	r.Handle("/v1/calendar/{id}", requireAuthenticatedUser(httpTransport.NewServer(
		ep.DeleteEventEndpoint,
		decodeHTTPDeleteEventRequest,
		encodeResponse))).Methods(http.MethodDelete)

	r.Handle("/v1/signup", httpTransport.NewServer(
		ep.SignUpEndpoint,
		decodeSignUpRequest,
		encodeResponse)).Methods(http.MethodPost)

	r.Handle("/v1/login", httpTransport.NewServer(
		ep.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse)).Methods(http.MethodPost)

	r.Handle("/v1/logout", httpTransport.NewServer(
		ep.LogoutEndpoint,
		decodeLogoutRequest,
		encodeResponse)).Methods(http.MethodPost)

	return r
}

func decodeHTTPListEventRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ListEventRequest

	req.UserId = context2.GetUser(r).UserID
	return req, nil
}

func decodeHTTPAddEventRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddEventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	req.Event.UserId = context2.GetUser(r).UserID
	return req, nil
}

func decodeHTTPDeleteEventRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DeleteEventRequest

	var err error
	req.EventId = mux.Vars(r)["id"]
	if err != nil {
		logger.Log("err", err)
		return nil, err
	}
	req.UserId = context2.GetUser(r).UserID
	return req, nil
}

func decodeSignUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LogoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch err {
	case repository.ErrRecordNotFound:
		w.WriteHeader(http.StatusNotFound)
	case repository.ErrDuplicateEmail:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"error": err})
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
