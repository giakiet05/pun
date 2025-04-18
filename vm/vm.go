package vm

import (
	"fmt"
	"pun/bytecode"
	"pun/error"
)

type Scope struct {
	Locals []interface{} // Biến local trong scope này
	Parent *Scope        // Scope cha (cho nested blocks)
}

type VM struct {
	Constants    []interface{}              // Pool hằng số (copy từ compiler)
	Code         []bytecode.Instruction     // Chương trình bytecode
	Stack        []interface{}              // Stack thực thi
	Globals      []interface{}              // Bộ nhớ global (tương ứng GlobalSymbol trong compiler)
	ScopeStack   []*Scope                   // Scope stack (lưu biến local)
	CurrentScope *Scope                     //Scope hiện tại
	Sp           int                        // Stack pointer
	Ip           int                        // Instruction pointer
	Builtins     map[string]BuiltinFunction //Lưu built-in function
	Errors       []customError.RuntimeError
}

func NewVM(constants []interface{}, code []bytecode.Instruction, globalsSize int) *VM {
	vm := &VM{
		Constants:  constants,
		Code:       code,
		Stack:      make([]interface{}, 0, 1024),
		Globals:    make([]interface{}, globalsSize),
		ScopeStack: make([]*Scope, 0),
		Sp:         -1,
		Ip:         0,
		Builtins:   make(map[string]BuiltinFunction),
	}

	// Khởi tạo global scope (root scope)
	globalScope := &Scope{
		Locals: make([]interface{}, 0), // Global scope không có locals
		Parent: nil,
	}
	vm.ScopeStack = append(vm.ScopeStack, globalScope)
	vm.CurrentScope = globalScope

	vm.Builtins["print"] = vm.builtinPrint
	vm.Builtins["ask"] = vm.builtinAsk

	return vm
}

func (v *VM) Run() {
	for v.Ip < len(v.Code) {

		//Nếu có lỗi thì in lỗi rồi return luôn
		if v.HasErrors() {
			return
		}

		inst := v.Code[v.Ip]

		switch inst.Op {
		case bytecode.OP_LOAD_CONST:
			val := v.Constants[inst.Operand.(int)]
			v.push(val)
		case bytecode.OP_LOAD_NOTHING:
			v.push(nil)
		case bytecode.OP_LOAD_GLOBAL:
			slot := inst.Operand.(int)
			if slot >= len(v.Globals) {
				v.addError(fmt.Sprintf("global variable slot %d out of bounds", slot), 0, 0, "runtime")
				continue
			}
			v.push(v.Globals[slot])

		case bytecode.OP_STORE_GLOBAL:
			slot := inst.Operand.(int)
			if slot >= len(v.Globals) {
				v.addError(fmt.Sprintf("global variable slot %d out of bounds", slot), 0, 0, "runtime")
				continue
			}

			v.Globals[slot] = v.pop()

		case bytecode.OP_LOAD_LOCAL:
			v.executeLoadLocal(inst.Operand.(*bytecode.LocalVar))

		case bytecode.OP_STORE_LOCAL:
			v.executeStoreLocal(inst.Operand.(*bytecode.LocalVar))

		case bytecode.OP_ENTER_SCOPE:
			localSize := inst.Operand.(int)
			v.pushScope(localSize)

		case bytecode.OP_LEAVE_SCOPE:
			v.popScope()

		case bytecode.OP_CALL:
			argCount := inst.Operand.(int)
			v.executeCall(argCount)
		case bytecode.OP_JUMP:
			offset := inst.Operand.(int)
			v.Ip += offset
		case bytecode.OP_JUMP_IF_FALSE:
			//Kiểm tra điều kiện (đã tính trước và lưu vào stack)
			//Nếu sai thì nhảy (tăng v.Ip)
			condition := v.pop().(bool)
			if !condition {
				v.Ip += inst.Operand.(int)
			}
		case bytecode.OP_RETURN:
			v.executeReturn()
		case bytecode.OP_MAKE_ARRAY:
			size := inst.Operand.(int)
			v.executeMakeArray(size)

		case bytecode.OP_ARRAY_GET:
			v.executeArrayGet()

		case bytecode.OP_ARRAY_SET:
			v.executeArraySet()
		case bytecode.OP_MAKE_FUNCTION:
			v.executeMakeFunction()
		case bytecode.OP_ADD:
			v.executeArithmetic("+")
		case bytecode.OP_SUB:
			v.executeArithmetic("-")
		case bytecode.OP_MUL:
			v.executeArithmetic("*")
		case bytecode.OP_DIV:
			v.executeArithmetic("/")
		case bytecode.OP_MOD:
			v.executeArithmetic("%")
		case bytecode.OP_POW:
			v.executeArithmetic("**")
		case bytecode.OP_EQ:
			v.executeComparison("==")
		case bytecode.OP_NEQ:
			v.executeComparison("!=")
		case bytecode.OP_GT:
			v.executeComparison(">")
		case bytecode.OP_GTE:
			v.executeComparison(">=")
		case bytecode.OP_LT:
			v.executeComparison("<")
		case bytecode.OP_LTE:
			v.executeComparison("<=")
		case bytecode.OP_AND:
			v.executeLogical("&&")
		case bytecode.OP_OR:
			v.executeLogical("||")
		case bytecode.OP_NOT:
			v.executeNot()
		case bytecode.OP_NEG:
			v.executeNegate()
		default:
			v.addError(fmt.Sprintf("unknown opcode: %s", inst.Op), 0, 0, "runtime")
		}
		v.Ip++
	}
}
