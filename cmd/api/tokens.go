package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mensurui/todoList/internal/data"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.logger.Fatal("error decoding")
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)

	if err != nil {
		app.logger.Printf("error while getting the user by the email: %s", err)
		return
	}

	match, err := user.Password.Match(input.Password)

	if err != nil {
		app.logger.Print("Error while matching the passwords")
		return
	}

	if !match {
		app.logger.Print("Password or email not correct")
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.logger.Printf("Error while generating token: %s", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, token, nil)
	if err != nil {
		app.logger.Println("Error while writing the token")
		return
	}

}
