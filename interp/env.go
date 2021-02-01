package interp

import "gitlab.com/rcmonnet/glox/lang"

// Env represents the interpreter state.
type env struct {
	values    map[string]interface{}
	enclosing *env
}

// NewEnv creates a new environment.
func newEnv(enclosing *env) *env {

	env := &env{}
	env.values = make(map[string]interface{})
	env.enclosing = enclosing
	return env
}

// define binds a variable name and its value for the environment.
// If the variable was already bound, the name value is bound
// instead and the old value is discarded.
func (e *env) define(name string, value interface{}) {

	e.values[name] = value
}

// get retrieves the value associated with a variable.
// It the variable is not bound a runtimeError is triggered.
func (e *env) get(name *lang.Token) interface{} {

	value, ok := e.values[name.Lexeme]
	if ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	panic(runtimeError{name,
		"Undefined variable '" + name.Lexeme + "'."})
}

// assign bind a new value with an existing variable.
// It returns a RuntimeError if the variable doesn't exist.
func (e *env) assign(name *lang.Token, value interface{}) {

	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}
	panic(runtimeError{name,
		"undefined variable '" + name.Lexeme + "'"})
}
