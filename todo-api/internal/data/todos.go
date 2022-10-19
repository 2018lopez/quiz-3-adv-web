//Filename: internal/data/todo.go

package data

import (
	"context"
	"database/sql"
	"errors"
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

// Get() allow us to retrieve a specific todo task by id
func (m TodoModel) Get(id int64) (*Todo, error) {

	//Ensure id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//query to get todo task by id
	query :=
		`
		SELECT id, created_at, title, description, completed FROM todo
		WHERE id = $1

	`

	//Declare Todo variable to hold return results

	var todo Todo

	//create context,
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup memory
	defer cancel()

	//Execute the query using QueryRow()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
	)
	if err != nil {
		//check type of err
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	//Success
	return &todo, nil
}

// Update() allow update todo task by id
func (m TodoModel) Update(todo *Todo) error {

	//query to update todo task record

	query :=
		`
		UPDATE todo 
		SET title = $1, description = $2, completed = $3
		WHERE id = $4
		RETURNING id
		
	`
	//create context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup to prevent memory leak
	defer cancel()

	args := []interface{}{
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.ID,
	}

	//check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}

	}

	return nil
}
