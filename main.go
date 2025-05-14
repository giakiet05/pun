package main

import (
	"fmt"
	"os"
	"pun/compiler"
	"pun/lexer"
	"pun/parser"
	"pun/repl"
	"pun/vm"
	"time"
)

func debug() {
	for { // Thêm vòng lặp để xử lý nhập sai
		var ch int
		fmt.Println("(1) Lexer (2) Parser (3) Compiler (4) VM (0) Exit")
		fmt.Print("Choose: ")
		_, err := fmt.Scanln(&ch)
		if err != nil {
			fmt.Println("Invalid input! Please enter a number.")
			continue
		}

		switch ch {
		case 1:
			repl.StartLexer()
			return // Thoát sau khi chạy xong
		case 2:
			repl.StartParser()
			return // Thoát sau khi chạy xong
		case 3:
			repl.StartCompiler()
			return
		case 4:
			repl.StartVM()
			return
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Unknown choice! Choose again.")
		}
	}
}

func run(filename ...string) {
	if len(filename) > 0 && filename[0] != "" {
		runFile(filename[0]) // Nếu có file, chạy file đó
		return
	}

	if len(os.Args) > 1 {
		runFile(os.Args[1]) // Run .pun file
		return
	}
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	l := lexer.NewLexer(string(data))
	p := parser.NewParser(l)
	c := compiler.NewCompiler()

	program := p.ParseProgram()

	if p.HasErrors() {
		p.PrintErrors()
		return
	}

	c.CompileProgram(program)

	if c.HasErrors() {
		c.PrintErrors()
		return
	}

	v := vm.NewVM(c.Constants, c.Code, len(c.GlobalSymbols))
	v.Run()

	if v.HasErrors() {
		v.PrintErrors()
		return
	}

}

func measureTime(fn func()) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("\nExecution time: %s\n", elapsed)
	}()
	fn()
}

func main() {
	measureTime(func() {
		run("example.pun")
		//debug()

	})
}
