// repl_compiler.go
package repl

import (
	"bufio"
	"fmt"
	"os"
	"pun/compiler"
	"pun/lexer"
	"pun/parser"
	"strings"
)

const COMPILER_PROMPT = "pun(compile)> "

func StartCompiler() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("🔥 Pun Compiler REPL - Bytecode Visualizer 🔥")
	fmt.Println("Input code to see compiled bytecode")

	for {
		fmt.Print(COMPILER_PROMPT)
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		l := lexer.NewLexer(input)
		p := parser.NewParser(l)
		c := compiler.NewCompiler()

		program := p.ParseProgram()

		if p.HasErrors() {
			p.PrintErrors()
			continue
		}

		c.CompileProgram(program)

		if c.HasErrors() {
			c.PrintErrors()
			continue
		}

		printCompilerState(c)
	}
}

func printCompilerState(c *compiler.Compiler) {
	// Print bytecode
	fmt.Println("\n=== BYTECODE ===")
	for i, ins := range c.Code {
		fmt.Printf("%3d: %d\n", i, ins)
	}

	// Print constants
	if len(c.Constants) > 0 {
		fmt.Println("\n=== CONSTANT POOL ===")
		for i, constVal := range c.Constants {
			fmt.Printf("%d: %#v\n", i, constVal)
		}
	}

	// Print globals
	if len(c.GlobalSymbols) > 0 {
		fmt.Println("\n=== GLOBAL SYMBOLS ===")
		for name, idx := range c.GlobalSymbols {
			fmt.Printf("GLOBAL[%d]: %s\n", idx, name)
		}
	}

	// Print scopes
	if len(c.Scopes) > 0 {
		fmt.Println("\n=== SCOPE STACK ===")
		for scopeLevel, scope := range c.Scopes {
			fmt.Printf("Scope %d:\n", scopeLevel)
			for name, idx := range scope {
				fmt.Printf("  LOCAL[%d]: %s\n", idx, name)
			}
		}
	}

	fmt.Println(strings.Repeat("═", 30))
}
