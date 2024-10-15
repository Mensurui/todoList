package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	route := httprouter.New()
	route.HandlerFunc(http.MethodGet, "/v1/healthchecker", app.healthChecker)
	route.HandlerFunc(http.MethodGet, "/v1/todo/:id", app.authentication(app.readTodoHandler))
	route.HandlerFunc(http.MethodPost, "/v1/todo", app.authentication(app.createTodoHandler))
	route.HandlerFunc(http.MethodPut, "/v1/todo/:id", app.authentication(app.updateTodoHandler))
	route.HandlerFunc(http.MethodGet, "/v1/todo", app.authentication(app.readAllTodoHandler))

	route.HandlerFunc(http.MethodPost, "/v1/user", app.userRegistrationHandler)
	route.HandlerFunc(http.MethodPost, "/v1/authentication", app.createAuthenticationTokenHandler)
	return route
}
