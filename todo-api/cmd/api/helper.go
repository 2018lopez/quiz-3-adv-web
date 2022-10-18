//filename: cmd/api/helper.go

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"todo.imerlopez.net/internal/validator"
)

//Define a new type

type envelope map[string]interface{}

func (app *application) readIdParam(r *http.Request) (int64, error) {

	//ParamsFromContext() function to get the request context as a slice
	params := httprouter.ParamsFromContext(r.Context())

	//get id from params
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {

		return 0, errors.New("invalid id parament")
	}

	return id, nil

}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	//convert map result into JSON data
	js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	//For newline,
	js = append(js, '\n')

	//add the headers

	for key, value := range headers {
		w.Header()[key] = value
	}

	//specify that well serve our response using json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	//write json as the http body

	w.Write(js)

	return nil

}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	//user http.maxByteReader() to limit the size of the request body to 1 mb 2 ^20
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)

	//check for bad request

	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		//switch to check for the error

		switch {
		//check for syntax error
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed json(at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("Body contains badly-formed JSON")
		//check for wrong types passed by the client
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contain incorrect JSON type for field %q", unmarshalTypeError)
			}
			return fmt.Errorf("body contains incorrect JSON type - at character %d", unmarshalTypeError.Offset)
		//Empty Body
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		//unmappable field
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "jsonL unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		//Large Files
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not larger than %d bytes", maxBytes)

		//Pass non-nil pointer error
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err

		}
	}

	//call decode again
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// the readString method returns a string value from the query parameters
// or returns a default value if no matching key is found
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	//get the value
	value := qs.Get(key)

	if value == "" {

		return defaultValue
	}

	return value
}

// the readCSV method split a value into a slice based on the comma character
// if no matching key is found then the default value is returned

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {

	//get value
	value := qs.Get(key)

	if value == "" {
		return defaultValue
	}

	//split the string base on the comma delimeter
	return strings.Split(value, ",")
}

// the read int method converts a string value to int value
// if the value cannot converted to an integer then a validation error is added to
// the validation errors map

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {

	//ge the value

	value := qs.Get(key)

	if value == "" {
		return defaultValue
	}

	//Perform the conversion of an int

	intValue, err := strconv.Atoi(value)

	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return intValue
}
