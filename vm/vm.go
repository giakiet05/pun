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

func (v *VM) pushScope(localSize int) {
	scope := &Scope{Locals: make([]interface{}, localSize), Parent: v.CurrentScope}
	v.ScopeStack = append(v.ScopeStack, scope)
	v.CurrentScope = scope
}

func (v *VM) popScope() {
	if len(v.ScopeStack) > 1 { // Giữ lại global scope
		v.ScopeStack = v.ScopeStack[:len(v.ScopeStack)-1]
		v.CurrentScope = v.ScopeStack[len(v.ScopeStack)-1]
	}
}

func (v *VM) Run() {
	for v.Ip < len(v.Code) {

		//Nếu có lỗi thì in lỗi rồi return luôn
		if v.HasErrors() {
			return
		}

		inst := v.Code[v.Ip]
		v.Ip++

		switch inst.Op {
		case bytecode.OP_LOAD_CONST:
			val := v.Constants[inst.Operand.(int)]
			v.push(val)

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
					}

					result := builtin(args...)
					if result != nil {
						v.push(result)
					}
				} else {
					v.addError("undefined builtin function", 0, 0, f)
				}
			default:
				v.addError("not callable", 0, 0, fmt.Sprintf("%T", fn))
			}

		case bytecode.OP_JUMP:
			v.Ip += inst.Operand.(int)
		case bytecode.OP_JUMP_IF_FALSE:
			//Kiểm tra điều kiện (đã tính trước và lưu vào stack)
			//Nếu sai thì nhảy (tăng v.Ip)
			condition := v.pop().(bool)
			if !condition {
				v.Ip += inst.Operand.(int)
			}

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
	}
}

func (v *VM) getScope(depth int) *Scope {
	scope := v.CurrentScope
	for i := 1; i < depth; i++ {
		scope = scope.Parent
	}
	return scope
}
