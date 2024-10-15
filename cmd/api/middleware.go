package main

import (
	"net/http"
	"strings"
)

func (app *application) authentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.logger.Printf("Error on header: %s", headerParts)
			return
		}

		token := headerParts[1]
		user, err := app.models.Users.GetByToken(token)

		if err != nil {
			app.logger.Printf("error getting the token: %s", err)
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
