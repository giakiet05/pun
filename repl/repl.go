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

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		pause()
		return
	}

	env := evaluator.NewEnvironment()
	runPunCode(string(data), env)

	pause() // Đợi user nhấn Enter trước khi thoát
}

func pause() {
	fmt.Print("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// runPunCode lexes, parses, and evaluates Pun code
func runPunCode(input string, env *evaluator.Environment) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	//If there are errors, then the evaluator will never be called
	if p.HasErrors() {
		p.PrintErrors()
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Printf("%v\n", evaluated)
	}
}
