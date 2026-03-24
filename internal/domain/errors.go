package domain

import "errors"

// ErrNotFound indicates the requested domain resource does not exist.
var ErrNotFound = errors.New("resource not found")
