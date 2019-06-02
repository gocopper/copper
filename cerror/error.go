// Package cerror provides a custom error type that can hold more context than the stdlib error package.
// The goal is to provide better logging and debugging.
package cerror

import (
	"fmt"
	"reflect"
	"strings"
)

// Error is a custom error type that can hold tags and cause of an error for better debugging and logging.
type Error struct {
	Message string
	Tags    map[string]interface{}
	Cause   error
}

// New is used to create a new Error with a cause, message, tags.
// Cause and tags are optional and can be nil.
func New(cause error, msg string, tags map[string]interface{}) error {
	return Error{
		Message: msg,
		Tags:    tags,
		Cause:   cause,
	}
}

// WithTags wraps the given error with tags
func WithTags(err error, tags map[string]interface{}) error {
	return Error{
		Message: err.Error(),
		Tags:    tags,
		Cause:   Cause(err),
	}
}

// Error creates a log-friendly string of the error using the cause and tags.
func (e Error) Error() string {
	var err strings.Builder
	var tags []string

	err.WriteString(e.Message)

	for tag, val := range e.Tags {
		if reflect.ValueOf(val).Kind() == reflect.Ptr && val.(*interface{}) != nil {
			tags = append(tags, fmt.Sprintf("%s=%+v", tag, *(val.(*interface{}))))
		} else {
			tags = append(tags, fmt.Sprintf("%s=%+v", tag, val))
		}
	}

	if len(tags) > 0 {
		err.WriteString(" where ")
		err.WriteString(strings.Join(tags, ","))
	}

	if e.Cause != nil {
		err.WriteString(" because\n")
		err.WriteString("> ")
		err.WriteString(e.Cause.Error())
	}

	return err.String()
}
