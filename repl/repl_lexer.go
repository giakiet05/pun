package repl

import (
	"bufio"
	"fmt"
	"os"
	"pun/lexer"
)

// StartLexer runs a REPL to test the lexer, and can also process a .pun file
func StartLexer() {
	if len(os.Args) > 1 {
		lexFile(os.Args[1])
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Pun Lexer REPL!")
	fmt.Println("Type your code below (type 'exit' to quit):")

	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			return
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		lexInput(input)
	}
}

// lexFile reads a .pun file and prints tokens
func lexFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	lexInput(string(data))
}

// lexInput tokenizes input and prints tokens
func lexInput(input string) {
	l := lexer.NewLexer(input)

	for {
		tok := l.NextToken()
		fmt.Printf("{Type:%s, Value:%q, Line:%d, Col:%d}\n", tok.Type, tok.Value, tok.Line, tok.Col)
		if tok.Type == lexer.TOKEN_EOF {
			break
		}
	}
}
