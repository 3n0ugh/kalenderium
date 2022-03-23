package errs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTP error responses

func errorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	js, err := json.Marshal(map[string]interface{}{"err": message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(js)
}

func ServerErrorResponse(w http.ResponseWriter) {
	message := "the server encountered a problem and could not process your request"
	errorResponse(w, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {

	message := "the requested resource could not be found"
	errorResponse(w, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, http.StatusMethodNotAllowed, message)
}

func RateLimitExceededResponse(w http.ResponseWriter) {
	message := "rate limit exceeded"
	errorResponse(w, http.StatusTooManyRequests, message)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	errorResponse(w, http.StatusUnauthorized, message)
}

func AuthenticationRequiredResponse(w http.ResponseWriter) {
	message := "you must be authenticated to access this resource"
	errorResponse(w, http.StatusUnauthorized, message)
}
