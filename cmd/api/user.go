package main

import (
	"encoding/json"
	"net/http"

	"github.com/Mensurui/todoList/internal/data"
)

func (app *application) userRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.logger.Println("Error Decoding")
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.logger.Println("Error setting the password")
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		app.logger.Println("Error inserting in the db")
		return
	}

	err = app.writeJSON(w, http.StatusOK, user, nil)

	if err != nil {
		app.logger.Println("Error writing the json")
	}
}
