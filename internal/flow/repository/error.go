package repository

import "errors"

var ErrNotFound = errors.New("not found")
var ErrTagExists = errors.New("tag exists")
var ErrEntityNotTagged = errors.New("entity not tagged")
var ErrPoolFull = errors.New("entity pool full")
var ErrWrongMRID = errors.New("mRID doesn't match")
var ErrBadData = errors.New("bad data")
