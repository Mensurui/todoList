package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	route := httprouter.New()
	route.HandlerFunc(http.MethodGet, "/v1/healthchecker", app.healthChecker)
	route.HandlerFunc(http.MethodGet, "/v1/todo/:id", app.readTodoHandler)
	route.HandlerFunc(http.MethodPost, "/v1/todo", app.createTodoHandler)

	return route
}
