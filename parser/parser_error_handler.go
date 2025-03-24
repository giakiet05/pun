package parser

import (
	"fmt"
	"pun/errors"
)

// addError records an error

func (p *Parser) addError(message string, line, col int) {
	p.errors = append(p.errors, errors.PunError{Message: message, Line: line, Column: col})
}

// HasErrors returns true if the parser has errors

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

// PrintErrors prints all errors

func (p *Parser) PrintErrors() {
	if p.HasErrors() {
		fmt.Println("Parsing errors:")
		for _, err := range p.errors {
			fmt.Println(err.Error())
		}
	}
}
