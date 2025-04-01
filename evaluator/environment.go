package evaluator

type Environment struct {
	store map[string]interface{}
	outer *Environment // Env cha (nếu có)
}

// Tạo môi trường mới (global)
func NewEnvironment() *Environment {
	env := &Environment{store: make(map[string]interface{}), outer: nil}

	env.Set("explode", NewExplodeFunction())

	return env
}

// Tạo môi trường con kế thừa từ env cha
func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{store: make(map[string]interface{}), outer: outer}
}

// Gán biến vào env hiện tại
func (env *Environment) Set(name string, value interface{}) {
	// Kiểm tra xem biến có tồn tại trong outer không
	if env.outer != nil {
		if _, exists := env.outer.Get(name); exists {
			env.outer.Set(name, value) // ✅ Cập nhật ở outer
			return
		}
	}

	env.store[name] = value // 🆕 Nếu chưa có trong outer, lưu vào env hiện tại
}

// Lấy giá trị biến (tìm trong env hiện tại, nếu không có thì tìm env cha)
func (env *Environment) Get(name string) (interface{}, bool) {
	val, ok := env.store[name]
	if !ok && env.outer != nil {
		return env.outer.Get(name) // Tìm tiếp trong env cha
	}
	return val, ok
}

// Kiểm tra giá trị có phải "truthy" không (dùng cho điều kiện)
func isTruthy(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		return v != ""
	default:
		return val != nil
	}
}
