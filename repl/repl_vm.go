// repl_vm.go
package repl

import (
	"bufio"
	"fmt"
	"os"
	"pun/compiler"
	"pun/lexer"
	"pun/parser"
	"pun/vm"
	"strings"
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

		machine := vm.NewVM(c.Constants, c.Code, len(c.GlobalSymbols))
		machine.Run()

		if machine.HasErrors() {
			machine.PrintErrors()
			continue
		}

		printVMState(machine)
	}
}

func printVMState(m *vm.VM) {
	fmt.Println("\n=== VM STATE ===")

	// Print IP and SP
	fmt.Printf("IP: %d\n", m.Ip)
	fmt.Printf("SP: %d\n", m.Sp)

	// Print stack
	if m.Sp >= 0 {
		fmt.Println("\n--- Stack ---")
		for i := m.Sp; i >= 0; i-- {
			fmt.Printf("[%d] %#v\n", i, m.Stack[i])
		}
	}

	// Print globals
	if len(m.Globals) > 0 {
		fmt.Println("\n--- Globals ---")
		for i, val := range m.Globals {
			if val != nil {
				fmt.Printf("GLOBAL[%d]: %#v\n", i, val)
			}
		}
	}

	fmt.Println(strings.Repeat("═", 30))
}
