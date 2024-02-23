package storage

import "errors"

var (
	ErrUrlNotFound = errors.New("url is not found")
	ErrUrlExists   = errors.New("url already exists")
)
