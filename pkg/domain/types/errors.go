package types

import "github.com/m-mizutani/goerr"

var (
	ErrInvalidOption = goerr.New("invalid option")
	ErrInvalidTask   = goerr.New("invalid task")

	ErrTestFailed = goerr.New("test fail")
)
