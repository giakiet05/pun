package vm

import (
	"fmt"
	"math"
	"pun/bytecode"
)

func (v *VM) executeArithmetic(op string) {
	if v.Sp < 1 {
		v.addError("stack underflow", 0, 0, "arithmetic operation")
		return
	}

	right := v.pop()
	left := v.pop()

	leftVal, ok1 := left.(float64)
	rightVal, ok2 := right.(float64)

	if !ok1 || !ok2 {
		v.addError(fmt.Sprintf("operations only supported between numbers, got %T and %T", left, right), 0, 0, "arithmetic operation")
		return
	}

	var result float64
	switch op {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	case "/":
		if rightVal == 0 {
			v.addError("division by zero", 0, 0, "arithmetic operation")
			return
		}
		result = leftVal / rightVal
	case "%":
		if rightVal == 0 {
			v.addError("division by zero", 0, 0, "arithmetic operation")
			return
		}
		result = float64(int64(leftVal) % int64(rightVal))
	case "**":
		result = math.Pow(leftVal, rightVal)
	default:
		v.addError(fmt.Sprintf("unsupported operator: %s", op), 0, 0, "arithmetic operation")
		return
	}

	v.push(result)
}

func (v *VM) executeComparison(op string) {
	if v.Sp < 1 {
		v.addError("stack underflow", 0, 0, "comparison operation")
		return
	}

	right := v.pop()
	left := v.pop()

	switch leftVal := left.(type) {
	case float64:
		rightVal, ok := right.(float64)
		if !ok {
			v.addError(fmt.Sprintf("cannot compare float with %T", right), 0, 0, "comparison operation")
			return
		}
		var result bool
		switch op {
		case "==":
			result = leftVal == rightVal
		case "!=":
			result = leftVal != rightVal
		case "<":
			result = leftVal < rightVal
		case ">":
			result = leftVal > rightVal
		case "<=":
			result = leftVal <= rightVal
		case ">=":
			result = leftVal >= rightVal
		default:
			v.addError(fmt.Sprintf("unsupported comparison operator: %s", op), 0, 0, "comparison operation")
			return
		}
		v.push(result)

	case string:
		rightVal, ok := right.(string)
		if !ok {
			v.addError(fmt.Sprintf("cannot compare string with %T", right), 0, 0, "comparison operation")
			return
		}
		switch op {
		case "==":
			v.push(leftVal == rightVal)
		case "!=":
			v.push(leftVal != rightVal)
		default:
			v.addError("string only supports == and != operators", 0, 0, "comparison operation")
			return
		}

	default:
		v.addError(fmt.Sprintf("unsupported type for comparison: %T", left), 0, 0, "comparison operation")
	}
}

func (v *VM) executeLogical(op string) {
	if v.Sp < 1 {
		v.addError("stack underflow", 0, 0, "logical operation")
		return
	}

	right := v.pop()
	left := v.pop()

	leftBool, ok1 := left.(bool)
	rightBool, ok2 := right.(bool)

	if !ok1 || !ok2 {
		v.addError("logical operators require boolean operands", 0, 0, "logical operation")
		return
	}

	var result bool
	switch op {
	case "&&":
		result = leftBool && rightBool
	case "||":
		result = leftBool || rightBool
	default:
		v.addError(fmt.Sprintf("unsupported logical operator: %s", op), 0, 0, "logical operation")
		return
	}

	v.push(result)
}

func (v *VM) executeNegate() {
	if v.Sp < 0 {
		v.addError("stack underflow", 0, 0, "unary operation")
		return
	}

	val := v.pop()
	if num, ok := val.(float64); ok {
		v.push(-num)
	} else {
		v.addError(fmt.Sprintf("cannot negate non-number type: %T", val), 0, 0, "unary operation")
	}
}

func (v *VM) executeNot() {
	if v.Sp < 0 {
		v.addError("stack underflow", 0, 0, "logical operation")
		return
	}

	val := v.pop()
	if b, ok := val.(bool); ok {
		v.push(!b)
	} else {
		v.addError(fmt.Sprintf("cannot logical NOT non-boolean type: %T", val), 0, 0, "logical operation")
	}
}

func (v *VM) executeLoadLocal(op *bytecode.LocalVar) {
	scope := v.getScope(op.Depth)

	if op.Slot >= len(v.CurrentScope.Locals) {
		v.addError(fmt.Sprintf("local variable slot %d out of bounds", op.Slot), 0, 0, "runtime")
		return
	}
	v.push(scope.Locals[op.Slot])
}

func (v *VM) executeStoreLocal(op *bytecode.LocalVar) {
	scope := v.getScope(op.Depth)

	if op.Slot >= len(v.CurrentScope.Locals) {
		v.addError(fmt.Sprintf("local variable slot %d out of bounds", op.Slot), 0, 0, "runtime")
		return
	}
	scope.Locals[op.Slot] = v.pop()
}

func (v *VM) executeCall(argCount int) {
	fn := v.pop()

	switch f := fn.(type) {
	case string: // Built-in function
		if builtin, ok := v.Builtins[f]; ok {
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = v.pop()
			}

			if argCount != len(args) {
				v.addError("wrong number of arguments for this function", 0, 0, f)
				return
			}

			result := builtin(args...)
			if result != nil {
				v.push(result)
			}
		} else {
			v.addError("undefined builtin function", 0, 0, f)
		}

	case *bytecode.Function: // User-defined function
		// Validate argument count
		if argCount != f.Arity {
			v.addError(fmt.Sprintf("expected %d arguments, got %d", f.Arity, argCount), 0, 0, f.Name)
			return
		}

		// Push a new scope for the function
		v.pushScope(f.LocalSize)

		// Set up local variables (parameters)
		for i := f.Arity - 1; i >= 0; i-- {
			v.CurrentScope.Locals[i] = v.pop()
		}

		// Save the current instruction pointer for returning
		v.push(v.Ip)

		// Jump to the function's start
		v.Ip = f.StartPC

	default:
		// Add more context to the error message
		v.addError(fmt.Sprintf("not callable: expected function, got %T (value: %v)", fn, fn), 0, 0, "execute call")
	}
}
func (v *VM) executeMakeArray(size int) {
	//Nếu số lượng phần tử trong array != op của make array thì lỗi
	if len(v.Stack) != size {
		v.addError("wrong array size", 0, 0, "make array")
		return
	}

	arr := make([]interface{}, size)
	//Thêm các phần từ trong stack vào arr
	for i := size - 1; i >= 0; i-- {
		arr[i] = v.pop()
	}

	v.push(arr)
}

func (v *VM) executeArrayGet() {
	indexInterface := v.pop() // Giả sử index luôn là int (nếu không, cần check thêm)
	arrInterface := v.pop()   // Lấy giá trị từ stack (kiểu interface{})

	indexFloat, ok := indexInterface.(float64)
	if !ok {
		v.addError(fmt.Sprintf("expected index to be a number, got %T instead", indexInterface), 0, 0, "array get")
		return
	}

	index := int(indexFloat)

	// Check 1: arr có phải slice không?
	arr, ok := arrInterface.([]interface{})
	if !ok {
		v.addError(fmt.Sprintf("expected array type, got %T", arrInterface), 0, 0, "array get")
		return
	}

	// Check 2: Index có hợp lệ không?
	if index < 0 || index >= len(arr) {
		v.addError(fmt.Sprintf("index %d out of bounds (array size: %d)", index, len(arr)), 0, 0, "array get")
		return
	}

	v.push(arr[index]) // Safe access!
}

func (v *VM) executeArraySet() {
	indexInterface := v.pop() // Giả sử index luôn là int (nếu không, cần check thêm)
	arrInterface := v.pop()   // Lấy giá trị từ stack (kiểu interface{})

	//Kiểm tra xem index có phải là float64 không (mặc định trong Pun kiểu number tương ứng với float64 trong Go)
	indexFloat, ok := indexInterface.(float64)
	if !ok {
		v.addError(fmt.Sprintf("expected index to be a number, got %T instead", indexInterface), 0, 0, "array get")
		return
	}
	//Sau đó chuyển thành int
	index := int(indexFloat)

	// Check 1: arr có phải slice không?
	arr, ok := arrInterface.([]interface{})
	if !ok {
		v.addError(fmt.Sprintf("expected array type, got %T", arrInterface), 0, 0, "array get")
		return
	}

	// Check 2: Index có hợp lệ không?
	if index < 0 || index >= len(arr) {
		v.addError(fmt.Sprintf("index %d out of bounds (array size: %d)", index, len(arr)), 0, 0, "array get")
		return
	}

	arr[index] = v.pop() //Lưu vào array
}

func (v *VM) executeMakeFunction() {
	// Pop the function object from the stack
	fnInterface := v.pop()

	// Ensure the popped value is of type *bytecode.Function
	fn, ok := fnInterface.(*bytecode.Function)
	if !ok {
		v.addError(fmt.Sprintf("expected function object, got %T", fnInterface), 0, 0, "make function")
		return
	}

	// Push the function object back onto the stack
	v.push(fn)
}

func (v *VM) executeReturn() {
	// 1. Pop the return value from the stack
	returnValue := v.pop()

	// 2. Restore the instruction pointer (IP) from the stack
	if v.Sp < 0 {
		v.addError("stack underflow: missing return address", 0, 0, "return")
		return
	}
	v.Ip = v.pop().(int)

	// 3. Exit the current scope
	v.popScope()

	// 4. Push the return value back onto the stack
	v.push(returnValue)
}
