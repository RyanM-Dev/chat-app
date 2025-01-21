package repositories

import "errors"

//general errors

var (
	ErrInvalidID       = errors.New("invalid ID format")
	ErrOperationFailed = errors.New("operation failed")
)

// chat errors
var (
	ErrChatNotFound         = errors.New("chat not found")
	ErrMissingChatParameter = errors.New("missing chat parameters")
	ErrDuplicateChat        = errors.New("duplicate chat")
)

//database errors

var (
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDuplicateKey       = errors.New("duplicate key error")
)
