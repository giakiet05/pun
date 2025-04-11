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
	fmt.Println("Input code to see compiled bytecode (JSON format)")

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

		c.CompileProgram(program)

		if p.HasErrors() {
			p.PrintErrors()
			continue
		}

		if c.HasErrors() {
			c.PrintErrors()
			continue
		}
		// In bytecode dạng JSON
		printBytecode(c)
	}
}

func printBytecode(c *compiler.Compiler) {
	// In bytecode dạng human-readable
	fmt.Println("\n=== BYTECODE ===")
	for i, inst := range c.Code {
		operand := ""
		if inst.Operand != nil {
			switch v := inst.Operand.(type) {
			case int:
				operand = fmt.Sprintf("%d", v)
			case bool:
				operand = fmt.Sprintf("%t", v)
			case string:
				operand = fmt.Sprintf("%q", v)
			default:
				operand = fmt.Sprintf("%v", v)
			}
		}
		fmt.Printf("%3d. %-15s %s\n", i, inst.Op, operand)
	}

	// In constant pool
	if len(c.Constants) > 0 {
		fmt.Println("\n=== CONSTANT POOL ===")
		for i, constVal := range c.Constants {
			fmt.Printf("%d: %#v\n", i, constVal)
		}
	}

	// In global symbols
	if len(c.GlobalSymbols) > 0 {
		fmt.Println("\n=== GLOBAL SYMBOLS ===")
		// Tạo slice sắp xếp theo index
		sortedGlobals := make([]string, len(c.GlobalSymbols))
		for name, idx := range c.GlobalSymbols {
			sortedGlobals[idx] = name
		}
		for idx, name := range sortedGlobals {
			fmt.Printf("GLOBAL[%d]: %s\n", idx, name)
		}
	}

	// In scope stack nếu có
	if len(c.Scopes) > 0 {
		fmt.Println("\n=== SCOPE STACK ===")
		for scopeLevel, scope := range c.Scopes {
			fmt.Printf("Scope %d:\n", scopeLevel)
			// Sắp xếp local vars theo index
			sortedLocals := make([]string, len(scope))
			for name, idx := range scope {
				sortedLocals[idx] = name
			}
			for idx, name := range sortedLocals {
				fmt.Printf("  LOCAL[%d]: %s\n", idx, name)
			}
		}
	}

	fmt.Println(strings.Repeat("═", 30))
}
