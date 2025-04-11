package repl

import (
	"bufio"
	"fmt"
	"os"
	"pun/compiler"
	"pun/lexer"
	"pun/parser"
	"pun/vm"
)

const VM_PROMPT = "pun(vm)> "

func StartVM() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("🔥 Pun VM REPL - Bytecode Executor 🔥")
	fmt.Println("Input code to execute (Press Ctrl+C to exit)")

	for {
		fmt.Print(VM_PROMPT)
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		// 1. Lexing
		l := lexer.NewLexer(input)

		// 2. Parsing
		p := parser.NewParser(l)
		program := p.ParseProgram()

		if p.HasErrors() {
			p.PrintErrors()
			continue
		}

		// 3. Compilation
		c := compiler.NewCompiler()
		c.CompileProgram(program)

		if c.HasErrors() {
			c.PrintErrors()
			continue
		}

		// 4. Execution
		machine := vm.NewVM(c.Constants, c.Code, len(c.GlobalSymbols))
		machine.Run() // Không cần check error vì đã xài addError()

		// 5. Check và hiển thị lỗi runtime nếu có
		if machine.HasErrors() {
			machine.PrintErrors()
			continue
		}

		// 6. Display results nếu không có lỗi
		printVMState(machine)
	}
}

func printVMState(m *vm.VM) {
	fmt.Println("\n=== EXECUTION RESULT ===")

	// Print global variables
	fmt.Println("\n--- Globals ---")
	for i, val := range m.Globals {
		if val != nil {
			fmt.Printf("$%d: %v\n", i, val)
		}
	}

	// Print stack (if not empty)
	if m.Sp >= 0 {
		fmt.Println("\n--- Stack ---")
		for i := 0; i <= m.Sp; i++ {
			fmt.Printf("[%d] %v\n", i, m.Stack[i])
		}
	}

	fmt.Println("═══")
}
