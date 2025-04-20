package repository

import "errors"

var ErrNotFound = errors.New("not found")
var ErrTagExists = errors.New("tag exists")
