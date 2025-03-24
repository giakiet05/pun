package errors

import "fmt"

// PunError is a custom error type for Pun language
type PunError struct {
	Message string
	Line    int
	Column  int
}

func (e *PunError) Error() string {
	return fmt.Sprintf("Error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}
