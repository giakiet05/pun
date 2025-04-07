package main

import (
	"fmt"
	"pun/cmd/repl"
)

func main() {
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
