package bytecode

type Function struct {
	Name      string
	Arity     int //số lượng param
	LocalSize int //Số lượng biến local (số lượng param + số lượng biến tạo trong hàm)
	StartPC   int //Địa chỉ bắt đầu thân hàm
}
