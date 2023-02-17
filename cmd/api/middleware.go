package api

// TODO: Rate limiting

import (
	"net/http"
	"strings"
)

func (a *App) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Authorization Header: <auth-scheme> <authorization-parameters>
		headerkv := strings.Split(authHeader, " ")
		if len(headerkv) != 2 || headerkv[0] != "apikey" {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}

		apikey := headerkv[1]

		ok, err := a.Auth.VerifyKey(apikey)
		if err != nil {
			a.serverErrorResponse(w, r, err)
			return
		}

		if !ok {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
