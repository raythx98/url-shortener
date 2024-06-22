package error_tool

import "fmt"

type AppError struct {
	code    int
	message string
	err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Error: %s, Code: %d, Err: %v", e.message, e.code, e.err)
}
