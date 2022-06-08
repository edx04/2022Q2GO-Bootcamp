package repository

import "errors"

var (
	ErrIdNotExist = errors.New("id does not exist")
	ErrOpenCsv    = errors.New("cannot open CSV file")
	ErrReadCsv    = errors.New("cannot read CSV file")
	ErrWriteCsv   = errors.New("cannot write CSV file")
)
