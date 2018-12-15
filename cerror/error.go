package cerror

import (
	"fmt"
	"strings"
)

type Error struct {
	Message string
	Tags    map[string]string
	Cause   error
}

func New(cause error, msg string, tags map[string]string) error {
	return Error{
		Message: msg,
		Tags:    tags,
		Cause:   cause,
	}
}

func (e Error) Error() string {
	var err strings.Builder
	var tags []string

	err.WriteString(e.Message)

	for tag, val := range e.Tags {
		tags = append(tags, fmt.Sprintf("%s=%s", tag, val))
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
