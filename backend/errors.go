package backend

import "errors"

var (
	// ErrTableAlreadyExists occures when creating table exists
	ErrTableAlreadyExists = errors.New("the table already exists")

	// ErrTableNotFound occures when creating table exists
	ErrTableNotFound = errors.New("there is no such table")

	// ErrIndexNotFound occurs when a table doesn't contain given column.
	ErrIndexNotFound = errors.New("there is no index corresponding column name")
)
