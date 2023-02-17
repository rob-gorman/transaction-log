package api

import (
	"errors"
	"fmt"
	"net/http"
)

// helper to log original server errors
func (a *App) logError(r *http.Request, err error) {
	msg := fmt.Sprintf(
		"%v ; request_method: %s ; request_url: %s",
		err,
		r.Method,
		r.URL.String(),
	)
	a.Log.Err(msg)
}

func payloadResponse(w http.ResponseWriter, status int, data *[]byte, headers http.Header) error {
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(*data)

	return nil
}

func (a *App) errorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	data := []byte(fmt.Sprintf("error: %s", err.Error()))

	err = payloadResponse(w, status, &data, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

func (a *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)

	err = errors.New("the server encountered a problem and could not process your request")
	a.errorResponse(w, r, http.StatusInternalServerError, err)
}

func (a *App) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	err := errors.New("the requested resource could not be found")
	a.errorResponse(w, r, http.StatusNotFound, err)
}

func (a *App) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("the %s method is not supported for this resource", r.Method)
	a.errorResponse(w, r, http.StatusMethodNotAllowed, err)
}

func (a *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.errorResponse(w, r, http.StatusBadRequest, err)
}

func (a *App) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "apikey")

	err := errors.New("invalid or missing authentication token")
	a.errorResponse(w, r, http.StatusUnauthorized, err)
}
