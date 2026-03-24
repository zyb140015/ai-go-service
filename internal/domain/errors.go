package domain

import "errors"

// ErrNotFound indicates the requested domain resource does not exist.
var ErrNotFound = errors.New("resource not found")

// ErrConflict indicates the requested write conflicts with an existing resource.
var ErrConflict = errors.New("resource conflict")
