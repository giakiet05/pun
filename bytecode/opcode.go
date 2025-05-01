package bytecode

type Opcode byte

const (
	OP_LOAD_CONST Opcode = iota
	OP_LOAD_NOTHING
	OP_LOAD_GLOBAL
	OP_STORE_GLOBAL
	OP_LOAD_LOCAL
	OP_STORE_LOCAL
	OP_ENTER_SCOPE
	OP_LEAVE_SCOPE
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_POW
	OP_EQ
	OP_NEQ
	OP_GTE
	OP_LTE
	OP_GT
	OP_LT
	OP_AND
	OP_OR
	OP_NOT
	OP_NEG
	OP_JUMP
	OP_JUMP_IF_FALSE
	OP_CALL
	OP_RETURN
	OP_MAKE_ARRAY
	OP_ARRAY_GET
	OP_ARRAY_SET
	OP_MAKE_FUNCTION
)

// Số byte operand ứng với mỗi opcode
var OperandWidths = map[Opcode]int{
	OP_LOAD_CONST:    1,
	OP_LOAD_NOTHING:  0,
	OP_LOAD_GLOBAL:   1,
	OP_STORE_GLOBAL:  1,
	OP_LOAD_LOCAL:    2,
	OP_STORE_LOCAL:   2,
	OP_ENTER_SCOPE:   0,
	OP_LEAVE_SCOPE:   0,
	OP_ADD:           0,
	OP_SUB:           0,
	OP_MUL:           0,
	OP_DIV:           0,
	OP_MOD:           0,
	OP_POW:           0,
	OP_EQ:            0,
	OP_NEQ:           0,
	OP_GTE:           0,
	OP_LTE:           0,
	OP_GT:            0,
	OP_LT:            0,
	OP_AND:           0,
	OP_OR:            0,
	OP_NOT:           0,
	OP_NEG:           0,
	OP_JUMP:          2,
	OP_JUMP_IF_FALSE: 2,
	OP_CALL:          1,
	OP_RETURN:        0,
	OP_MAKE_ARRAY:    1,
	OP_ARRAY_GET:     0,
	OP_ARRAY_SET:     0,
	OP_MAKE_FUNCTION: 1,
}

// Encode opcode + operands thành []byte
func Make(op Opcode, operands ...int) []byte {
	ins := []byte{byte(op)}
	width := OperandWidths[op]

	for _, o := range operands {
		if width == 1 {
			ins = append(ins, byte(o))
		} else if width == 2 {
			ins = append(ins, byte(o>>8), byte(o))
		}
	}
	return ins
}

// Decode operands từ []byte
func ReadOperands(op Opcode, ins []byte) ([]int, int) {
	width := OperandWidths[op]
	offset := 0
	operands := []int{}

	if width == 1 {
		operands = append(operands, int(ins[0]))
		offset = 1
	} else if width == 2 {
		operands = append(operands, int(ins[0])<<8|int(ins[1]))
		offset = 2
	}
	return operands, offset
}
