package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mensurui/todoList/internal/data"
)

func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	//app.models.Create/title, description/

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.logger.Println("error in decoding")
	}

	todo := &data.Todo{
		Title:       input.Title,
		Description: input.Description,
	}

	err = app.models.Todos.Create(todo)

	if err != nil {
		app.logger.Println("Error when creating")
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/todo/%d", todo.ID))

	err = app.writeJSON(w, http.StatusOK, todo, headers)

	if err != nil {
		app.logger.Println("Error Writing the Response")
		return
	}
}

func (app *application) readTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readID(r)
	if err != nil {
		app.logger.Println("Error")
	}

	todo, err := app.models.Todos.Get(id)

	if err != nil {
		app.logger.Println("Error fetching from database")
	}
	err = app.writeJSON(w, http.StatusOK, todo, nil)
	if err != nil {
		app.logger.Println("Error")
	}

}
