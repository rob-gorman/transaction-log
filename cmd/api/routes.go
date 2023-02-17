package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *App) registerRoutes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v0/register", a.register)

	router.HandlerFunc(http.MethodPost, "/v0/logs", a.authenticate(a.createEvent))
	
	router.HandlerFunc(http.MethodGet, "/v0/logs/:field/:value", a.authenticate(a.getLogsByField))
	
	return router
}
