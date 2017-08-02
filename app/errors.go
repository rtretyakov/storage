package app

import "errors"

var (
	errNotFound  = errors.New("not found")
	errWrongType = errors.New("wrong type")
)
