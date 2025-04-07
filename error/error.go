package customError

import "fmt"

// Base struct cho mọi lỗi trong Pun
type PunError struct {
	Message string
	Line    int
	Column  int
}

// Implement error interface
func (e *PunError) Error() string {
	return fmt.Sprintf("Error at line (%d:%d): %s", e.Line, e.Column, e.Message)
}

// SyntaxError - Dành riêng cho parser
type SyntaxError struct {
	PunError        // Embedded struct
	Context  string // Thêm context code
}

// Override Error() cho SyntaxError
func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError at line (%d:%d): %s\nContext: %s",
		e.Line, e.Column, e.Message, e.Context)
}

type CompilationError struct {
	PunError
	Context string
}

func (e *CompilationError) Error() string {
	return fmt.Sprintf("CompilationError at line (%d:%d): %s\nContext: %s",
		e.Line, e.Column, e.Message, e.Context)
}

type RuntimeError struct {
	PunError
	Context string // Thông tin bổ sung
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError at line (%d:%d): %s\nContext: %s",
		e.Line, e.Column, e.Message, e.Context)
}
