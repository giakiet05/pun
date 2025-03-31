package evaluator

type StopException struct{}
type ContinueException struct{}
type ReturnException struct {
	Value interface{}
}
