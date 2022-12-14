//Filename: internal/data/todo.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// Delete() remove a todo task by id
func (m TodoModel) Delete(id int64) error {

	//Verify if id is valid
	if id < 1 {
		return ErrRecordNotFound
	}

	//Delete query
	query :=
		`
		DELETE FROM todo WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup to prevent memory leak
	defer cancel()

	//Execute Delete query
	result, err := m.DB.ExecContext(ctx, query, id)

	//Check for error
	if err != nil {
		return err

	}

	//check for how many rows were affected by the delete operation
	//call the rowAffected method on the result var

	rowsAffected, err := result.RowsAffected()

	if err != nil {

		return err
	}

	//check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

//get all method returns a list of all schools sort by id

func (m TodoModel) GetAll(title string, filters Filters) ([]*Todo, Metadata, error) {
	//construct query

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, title, description, completed
		FROM todo
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		
		ORDER BY %s %s, id ASC LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortOrder())
	//CREATE a 3 sec timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{title, filters.limit(), filters.offset()}
	//execute
	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	//close the result set
	defer rows.Close()

	//totalRecord
	totalRecords := 0

	//Initialize an empty slice to hold Todo data
	todos := []*Todo{}

	//iterate over the rows in the result set

	for rows.Next() {
		var todo Todo
		//scan the values from row into todo struct
		err := rows.Scan(
			&totalRecords,
			&todo.ID,
			&todo.CreatedAt,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		//add the todo to our slice
		todos = append(todos, &todo)

	}

	//Check for errors after looping through the result set

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// return slice of todos
	return todos, metadata, nil

}
