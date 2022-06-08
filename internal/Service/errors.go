package service

import "errors"

var (
	ErrIdNotExist     = errors.New("id does not exist")
	ErrOpenCsv        = errors.New("cannot open CSV file")
	ErrReadCsv        = errors.New("cannot read CSV file")
	ErrWriteCsv       = errors.New("cannot write CSV file")
	ErrParamItems     = errors.New("param items must be integer")
	ErrItemsPerWorker = errors.New("param items_per_worker must be integer")
)
