//Filename: internal/data/todo.go

package data

import (
	"context"
	"database/sql"
	"time"

	"todo.imerlopez.net/internal/validator"
)

type Todo struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
}

func ValidateTodo(v *validator.Validator, todo *Todo) {

	//use the check method to excute our validation check

	//Check for title if empty and size
	v.Check(todo.Title != "", "title", "must be provided")
	v.Check(len(todo.Title) <= 20, "title", "must not be more than 20 bytes long")

	//check for descriptions if empty and size
	v.Check(todo.Description != "", "description", "must be provided")
	// v.Check(len(todo.Descriptions) >= 8, "descriptions", "must be atleast 8 bytes")
}

//Define a TodoModel which wrap a sql.DB connection pool

type TodoModel struct {
	DB *sql.DB
}

// insert() create todo task
func (m TodoModel) Insert(todo *Todo) error {

	//insert query to add data to todo table

	query :=
		`	
		INSERT INTO todo(title, description, completed) 
		values($1,$2,$3)
		RETURNING id, created_at
	`
	args := []interface{}{todo.Title, todo.Description, todo.Completed}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup to prevent memory leak
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt)
}
