package evaluator

import "fmt"

func NewShoutFunction() *BuiltInFunction {
	return &BuiltInFunction{
		Fn: func(args ...interface{}) interface{} {
			for _, arg := range args {
				fmt.Print(arg, " ")
			}
			fmt.Println()
			return nil
		},
	}
}

func NewExplodeFunction() *BuiltInFunction {
	return &BuiltInFunction{
		Fn: func(args ...interface{}) interface{} {
			fmt.Println("BOOM!")
			return nil
		},
	}
}

func NewLenFunction() *BuiltInFunction {
	return &BuiltInFunction{
		Fn: func(args ...interface{}) interface{} {
			if len(args) != 1 {
				return fmt.Errorf("len() requires exactly one argument")
			}
			switch v := args[0].(type) {
			case string:
				return len(v)
			case []interface{}:
				return len(v)
			default:
				return fmt.Errorf("argument to len() must be a string or an array")
			}
		},
	}
}
