package pkg

import "errors"

var (
	ErrMethodNotAllowed = errors.New("method Not Allowed")
	FailedToStartServer = errors.New("failed to start server")
	ErrInternalError    = errors.New("internal error")
)
