package chttp

import (
	"errors"
	"io/fs"
)

// EmptyFS is a simple implementation of fs.FS interface that only returns an error.
// This implementation emulates an empty directory.
type EmptyFS struct{}

// Open returns an error since this fs is empty.
func (fs *EmptyFS) Open(string) (fs.File, error) {
	return nil, errors.New("empty fs")
}
