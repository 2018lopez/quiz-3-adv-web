//Filename: cmd/api/todo.go

package main

import (
	"errors"
	"fmt"
	"net/http"

	"todo.imerlopez.net/internal/data"
	"todo.imerlopez.net/internal/validator"
)

//Todo handler to create todo task - POST

func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {

	//our target decode destination

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	//initialize the new json.decoder instance

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//copy the values from the input struct to a new todo struct
	todo := &data.Todo{
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
	}

	//initialize a new validator instance

	v := validator.New()

	//check the map to determine if there were any validation errors
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//create todo task
	err = app.models.Todos.Insert(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}

	//create a location header for newly created resource: todo task
	headers := make(http.Header)
	headers.Set("Locations", fmt.Sprintf("/v1/todos/%d", todo.ID))

	//write json response
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// get todo task by id
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {

	//read id parameter
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//fetch the specifc todo tasks
	todo, err := app.models.Todos.Get(id)

	//handler errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	//write json data return by get
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, nil)

	}

}

// todo task update handler
func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {

	//get id from parameter
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//fetch record from the db
	todo, err := app.models.Todos.Get(id)

	//handler errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			fmt.Println("hd")
			app.notFoundResponse(w, r)
		default:
			fmt.Println("h")
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	//create an input struct to hold data read in from client
	//update input struct by pointer

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Completed   *bool   `json:"completed"`
	}

	//initialize the new json.decoder instance
	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	//check for updates

	if input.Title != nil {
		todo.Title = *input.Title
	}

	if input.Description != nil {
		todo.Description = *input.Description
	}

	if input.Completed != nil {
		todo.Completed = *input.Completed
	}

	//initialize a new Validator instance
	v := validator.New()

	//check the map to determine if there were any validation errors
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Pass the updated todo task record to update method
	err = app.models.Todos.Update(todo)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	//write data by get

	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)

	if err != nil {

		app.serverErrorResponse(w, r, err)

	}

}

//Delete handler: Delete todo task by id

func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {

	//get id
	id, err := app.readIdParam(r)

	if err != nil {

		app.notFoundResponse(w, r)
		return
	}

	//delete a todo task from the database. Send 404 not found status to client
	//if no matching record

	err = app.models.Todos.Delete(id)

	//Handler error
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	//Return 200 status ok to client if record is delete successfull
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Todo Task SuccessFully Deleted"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
