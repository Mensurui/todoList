package main

import (
	"net/http"
)

func (app *application) healthChecker(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, app.config.port, nil)
	if err != nil {
		app.logger.Fatal(err)
	}

}
