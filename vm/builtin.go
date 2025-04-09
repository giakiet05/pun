package vm

import (
	"bufio"
	"fmt"
	"os"
)

type BuiltinFunction func(args ...interface{}) interface{}

func (v *VM) builtinPrint(args ...interface{}) interface{} {
	fmt.Println(args...)
	return nil
}

func (v *VM) builtinAsk(args ...interface{}) interface{} {
	prompt := args[0].(string)
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func (v *VM) builtinLen(arg interface{}) int {
	return 1
}
