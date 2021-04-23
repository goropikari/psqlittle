package core

import "errors"

var (
	// ErrTableAlreadyExists occures when creating table exists
	ErrTableAlreadyExists = errors.New("the table already exists")

	// ErrIndexNotFound occurs when a table doesn't contain given column.
	ErrIndexNotFound = errors.New("there is no index corresponding column name")
)
