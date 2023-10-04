package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

/**
	Form struct which anonymously embeds a url.Values object
	(to hold the form data) and an Errors field to hold any validation
	errors for the form data.
**/
type Form struct {
	url.Values
	Errors errors
}

/**
	New function to intialize new Form struct, with incoming
	form data parameter
**/
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}


/**
	Validated Required fields exist
**/
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

/**
	Validates maximum length that a specific field in the form
	contains a maximum number of characters
**/
func (f *Form) MaxLength(field string, length int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > length {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (max is %d characters)", length))
	}
}


/**
	PermittedValues method to check that a specified field in the form
	matches one of a set of specific permitted values. If the check fails
	then add the appropriate message to the form errors
**/
func (f *Form) PermittedValues(field string, options ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, option := range options {
		if value == option {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}


/**
	Valid method which returns true if there are no errors
**/
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}