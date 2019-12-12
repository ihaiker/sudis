package daemon

import "errors"

var (
	ErrNotFound = errors.New("the program not found.")
	ErrExists   = errors.New("the program is exists.")
)
