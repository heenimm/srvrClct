package pkg

import "errors"

var (
	ErrMethodNotAllowed      = errors.New("method Not Allowed")
	FailedToStartServer      = errors.New("failed to start server")
	ErrInternalError         = errors.New("internal error")
	ErrInvalidExpression     = errors.New("invalid expression")
	ErrInsufficientOperands  = errors.New("insufficient operands")
	ErrMismatchedParentheses = errors.New("mismatched parentheses")
	ErrInvalidNumber         = errors.New("invalid number")
	ErrDivisionByZero        = errors.New("division by zero")
	ErrInvalidOperation      = errors.New("invalid operation")
)
