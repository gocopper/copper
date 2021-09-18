// Package cerrors provides a custom error type that can hold more context than the stdlib error package.
package cerrors

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Error can wrap an error with additional context such as structured tags.
type Error struct {
	Message string
	Tags    map[string]interface{}
	Cause   error
}

// New creates an error by (optionally) wrapping an existing error and
// annotating the error with structured tags.
func New(cause error, msg string, tags map[string]interface{}) error {
	return Error{
		Message: msg,
		Tags:    tags,
		Cause:   cause,
	}
}

// WithTags annotates an existing error with structured tags.
func WithTags(err error, tags map[string]interface{}) error {
	cerr, ok := err.(Error) //nolint:errorlint
	if !ok {
		return Error{
			Message: err.Error(),
			Tags:    tags,
			Cause:   errors.Unwrap(err),
		}
	}

	return Error{
		Message: cerr.Message,
		Tags:    tags,
		Cause:   cerr.Cause,
	}
}

// Unwrap returns the underlying cause of an error (if any).
func (e Error) Unwrap() error {
	return e.Cause
}

// Error returns a human-friendly string that contains the
// entire error chain along with all of the tags on each
// error.
func (e Error) Error() string {
	var (
		err  strings.Builder
		tags []string
	)

	err.WriteString(e.Message)

	for tag, val := range e.Tags {
		reflectVal := reflect.ValueOf(val)
		if reflectVal.Kind() == reflect.Ptr && !reflectVal.IsNil() {
			tags = append(tags, fmt.Sprintf("%s=%+v", tag, reflectVal.Elem()))
		} else {
			tags = append(tags, fmt.Sprintf("%s=%+v", tag, val))
		}
	}

	sort.Strings(tags)

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
