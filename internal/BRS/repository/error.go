package repository

import "errors"

var (
	ErrNotFound = errors.New("entity not found")
	ErrTagExists = errors.New("tag already exists")
	ErrPoolFull = errors.New("entity pool is full")
)
