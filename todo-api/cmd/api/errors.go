//Filenname: cmd/api/errors.go

package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

//Send JSON Format error message

func (app *application) errorRepsonse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	//create json response
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//server error responses

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log the error
	app.logError(r, err)

	//Prepare msg with the error
	message := "the server encountered a problem and couldn't process the request"
	app.errorRepsonse(w, r, http.StatusInternalServerError, message)
}

// The not found response

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	//create msg
	message := "Request resource couldn't be found"
	app.errorRepsonse(w, r, http.StatusNotFound, message)
}

// a method not allowed response

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	//create msg
	message := fmt.Sprintf("the %s method is not supported for this resources", r.Method)
	app.errorRepsonse(w, r, http.StatusMethodNotAllowed, message)
}

// bad request
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.errorRepsonse(w, r, http.StatusBadRequest, err.Error())

}

//Validation errors

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorRepsonse(w, r, http.StatusUnprocessableEntity, errors)
}

//edit error

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	//create msg
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorRepsonse(w, r, http.StatusConflict, message)
}

// Rate Limit Errors
func (app *application) rateLimitExceedeResponse(w http.ResponseWriter, r *http.Request) {
	//create msg
	message := "rate limit exceeded"
	app.errorRepsonse(w, r, http.StatusTooManyRequests, message)
}
