package api

import (
	"auditlog/utils"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	res, err := a.Auth.Register()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	fmt.Printf("register handler object: %s", res)

	payloadResponse(w, http.StatusOK, &res, r.Header)
}

func (a *App) createEvent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.badRequestResponse(w, r, fmt.Errorf("malformed request body"))
		return
	}

	a.Log.Info("object recieved: %s", string(req))

	err = a.DB.InsertEvent(utils.PrepareRequest(req))
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *App) getLogsByField(w http.ResponseWriter, r *http.Request) {
	field, err := utils.ReadFieldParam(r)
	if err != nil {
		a.badRequestResponse(w, r, err)
	}

	value, err := utils.ReadValueParam(r)
	if err != nil {
		a.badRequestResponse(w, r, err)
	}

	res, err := a.DB.SelectRowsByField(field, value)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	if len(res) == 0 || res == nil {
		a.notFoundResponse(w, r)
		return
	}

	payloadResponse(w, http.StatusOK, &res, r.Header)
}
