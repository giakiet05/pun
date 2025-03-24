package repl

import (
	"bufio"
	"fmt"
	"os"
	"pun/evaluator"
	"pun/lexer"
	"pun/parser"
)

const PROMPT = "pun> "

func Start() {
	if len(os.Args) > 1 {
		runFile(os.Args[1]) // Run .pun file
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	env := evaluator.NewEnvironment()

	fmt.Println("Welcome to the Pun language REPL!")
	fmt.Println("Type your Pun code below:")

	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			return
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		runPunCode(input, env)
	}
}

// runFile executes a .pun file
func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	env := evaluator.NewEnvironment()
	runPunCode(string(data), env)
}

// runPunCode lexes, parses, and evaluates Pun code
func runPunCode(input string, env *evaluator.Environment) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if p.HasErrors() {
		p.PrintErrors()
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Printf("%v\n", evaluated)
	}
}

// StartLexerREPL runs a REPL to test the lexer, and can also process a .pun file
func StartLexerREPL() {
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
