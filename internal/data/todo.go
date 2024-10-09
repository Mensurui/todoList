package data

import (
	"database/sql"
)

type Todo struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoModel struct {
	DB *sql.DB
}

func (t *TodoModel) Create(todo *Todo) error {
	query := `
	INSERT INTO todos(title, description)
	VALUES($1, $2)
	RETURNING id, title, description`

	args := []interface{}{todo.Title, todo.Description}
	err := t.DB.QueryRow(query, args...).Scan(&todo.ID, &todo.Title, &todo.Description)

	if err != nil {
		return err
	}

	return nil
}

func (t *TodoModel) Get(id int64) (*Todo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, title, description
	FROM todos
	WHERE id=$1`

	var todo Todo

	err := t.DB.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}
