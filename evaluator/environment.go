package evaluator

type Environment struct {
	store map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]interface{})}
}

func (env *Environment) Set(name string, value interface{}) {
	env.store[name] = value
}

func (env *Environment) Get(name string) (interface{}, bool) {
	val, ok := env.store[name]
	return val, ok
}
