package bytecode

type LocalVar struct {
	Slot  int
	Depth int
}

// Opcodes dạng UPPER_CASE y chang
const (
	OP_LOAD_CONST    = "LOAD_CONST"
	OP_LOAD_NOTHING  = "LOAD_NOTHING"
	OP_LOAD_GLOBAL   = "LOAD_GLOBAL"
	OP_STORE_GLOBAL  = "STORE_GLOBAL"
	OP_LOAD_LOCAL    = "LOAD_LOCAL"
	OP_STORE_LOCAL   = "STORE_LOCAL"
	OP_ENTER_SCOPE   = "ENTER_SCOPE"
	OP_LEAVE_SCOPE   = "LEAVE_SCOPE"
	OP_ADD           = "ADD"
	OP_SUB           = "SUB"
	OP_MUL           = "MUL"
	OP_DIV           = "DIV"
	OP_MOD           = "MOD"
	OP_POW           = "POW"
	OP_EQ            = "EQ"
	OP_NEQ           = "NEQ"
	OP_GTE           = "GTE"
	OP_LTE           = "LTE"
	OP_GT            = "GT"
	OP_LT            = "LT"
	OP_AND           = "AND"
	OP_OR            = "OR"
	OP_NOT           = "NOT"
	OP_NEG           = "NEG"
	OP_JUMP          = "JUMP"
	OP_JUMP_IF_FALSE = "JUMP_IF_FALSE"
	OP_CALL          = "CALL"
	OP_RETURN        = "RETURN"
	OP_MAKE_ARRAY    = "MAKE_ARRAY"
	OP_ARRAY_GET     = "ARRAY_GET"
	OP_ARRAY_SET     = "ARRAY_SET"
	OP_MAKE_FUNCTION = "MAKE_FUNCTION"
)

// Instruction struct giữ nguyên
type Instruction struct {
	Op      string      // Opcode dạng string
	Operand interface{} // Tham số
	Line    int         // Số dòng source
}
