// Filename: internal/validator/validator.go

package validator

import (
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	PhoneRX = regexp.MustCompile(`^\+?\(?[0-9]{3}\)?\s?-\s?[0-9]{3}\s?-\s?[0-9]{4}$`)
)

// we create a type that wraps our validation errors map
type Validator struct {
	Errors map[string]string
}

//New create a new validator instance

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

//valid() check the errors map for entries

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

//In () check if an element can be found in provide list of

func In(element string, list ...string) bool {
	for i := range list {
		if element == list[i] {
			return true
		}
	}
	return false
}

// matches() return true if a string value matches a specific regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

//addError () add an error entry to Error Map

func (v *Validator) AddError(key, message string) {

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}

}

//check() performs the validation checks and call the addError
//method in turn if an error entry needs to be added

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {

		v.AddError(key, message)
	}

}

//Unique() check that there are no repeating values in slices

func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
