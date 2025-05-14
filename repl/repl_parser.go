package repl

import (
	"bufio"
	"fmt"
	"os"
	"pun/ast"
	"pun/lexer"
	"pun/parser"
	"strings"
)

const PROMPT = "pun> "

func StartParser() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to Pun REPL! (Press Ctrl+C to exit)")
	fmt.Println("AST will be printed after each parse.")

	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			break // Thoát nếu gặp lỗi đọc input (vd: Ctrl+D)
		}

		input := scanner.Text()
		if input == "" {
			continue // Bỏ qua input trống
		}

		l := lexer.NewLexer(input)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		if p.HasErrors() {
			p.PrintErrors()
			continue // 🚨 Đừng exit, cho nhập lại!
		}

		// In AST đẹp hơn thay vì fmt.Println(program)
		printAST(program)
	}
}

func printAST(program *ast.Program) {
	fmt.Println("=== AST DUMP ===")
	for i, stmt := range program.Statements {
		fmt.Printf("[%d] %s\n", i+1, astToString(stmt))
	}
}

// Hàm chuyển AST node thành string dễ đọc
func astToString(node ast.Node) string {
	switch n := node.(type) {
	// ========== Statements ==========
	case *ast.AssignStatement:
		return fmt.Sprintf("ASSIGN: %s = %s",
			astToString(n.Name),
			astToString(n.Value))

	case *ast.BlockStatement:
		lines := []string{"BLOCK {"}
		for _, stmt := range n.Statements {
			lines = append(lines, "  "+astToString(stmt))
		}
		lines = append(lines, "}")
		return strings.Join(lines, "\n")

	case *ast.IfStatement:
		s := fmt.Sprintf("IF (%s) %s",
			astToString(n.Condition),
			astToString(n.Body))

		for _, elif := range n.ElseIfs {
			s += fmt.Sprintf(" ELIF (%s) %s",
				astToString(elif.Condition),
				astToString(elif.Body))
		}

		if n.ElseBlock != nil {
			s += fmt.Sprintf(" ELSE %s", astToString(n.ElseBlock.Body))
		}
		return s

	case *ast.ForStatement:
		return fmt.Sprintf("FOR (%s; %s; %s) %s",
			astToString(n.Init),
			astToString(n.Condition),
			astToString(n.Update),
			astToString(n.Body))

	case *ast.WhileStatement:
		return fmt.Sprintf("WHILE (%s) %s",
			astToString(n.Condition),
			astToString(n.Body))

	case *ast.UntilStatement:
		return fmt.Sprintf("UNTIL (%s) %s",
			astToString(n.Condition),
			astToString(n.Body))

	case *ast.ExpressionStatement:
		return astToString(n.Expression)

	case *ast.FunctionDefinitionStatement:
		params := []string{}
		for _, p := range n.Parameters {
			params = append(params, astToString(p))
		}
		return fmt.Sprintf("FUNC %s(%s) %s",
			astToString(n.Name),
			strings.Join(params, ", "),
			astToString(n.Body))

	case *ast.MethodDefinitionStatement:
		params := []string{}
		for _, p := range n.Parameters {
			params = append(params, astToString(p))
		}
		return fmt.Sprintf("METHOD %s.%s(%s) %s",
			astToString(n.Receiver),
			astToString(n.Name),
			strings.Join(params, ", "),
			astToString(n.Body))

	case *ast.BreakStatement:
		return "BREAK"

	case *ast.ContinueStatement:
		return "CONTINUE"

	case *ast.ReturnStatement:
		if n.Value != nil {
			return fmt.Sprintf("RETURN %s", astToString(n.Value))
		}
		return "RETURN"

	// ========== Expressions ==========
	case *ast.Identifier:
		return fmt.Sprintf("ID(%s)", n.Value)

	case *ast.NumberExpression:
		return fmt.Sprintf("NUM(%v)", n.Value)

	case *ast.StringExpression:
		return fmt.Sprintf("STR(%q)", n.Value)

	case *ast.BooleanExpression:
		return fmt.Sprintf("BOOL(%v)", n.Value)

	case *ast.NothingExpression:
		return "NOTHING"

	case *ast.UnaryExpression:
		return fmt.Sprintf("UNARY_OP(%s%s)",
			n.Operator,
			astToString(n.Value))

	case *ast.BinaryExpression:
		return fmt.Sprintf("BINARY_OP(%s %s %s)",
			astToString(n.Left),
			n.Operator,
			astToString(n.Right))

	case *ast.IncDecExpression:
		if n.IsPrefix {
			return fmt.Sprintf("INC_DEC(%s%s)",
				n.Operator,
				astToString(n.Value))
		}
		return fmt.Sprintf("INC_DEC(%s%s)",
			astToString(n.Value),
			n.Operator)

	case *ast.ArrayExpression:
		elements := []string{}
		for _, el := range n.Elements {
			elements = append(elements, astToString(el))
		}
		return fmt.Sprintf("ARRAY[%s]", strings.Join(elements, ", "))

	case *ast.ArrayIndexExpression:
		return fmt.Sprintf("INDEX(%s[%s])",
			astToString(n.Array),
			astToString(n.Index))

	case *ast.FunctionCallExpression:
		args := []string{}
		for _, arg := range n.Arguments {
			args = append(args, astToString(arg))
		}
		return fmt.Sprintf("CALL %s(%s)",
			astToString(n.Function),
			strings.Join(args, ", "))

	case *ast.MethodCallExpression:
		args := []string{}
		for _, arg := range n.Arguments {
			args = append(args, astToString(arg))
		}
		return fmt.Sprintf("METHOD_CALL %s.%s(%s)",
			astToString(n.Caller),
			n.Method,
			strings.Join(args, ", "))

	default:
		return fmt.Sprintf("UNKNOWN_NODE(%T)", n)
	}
}
