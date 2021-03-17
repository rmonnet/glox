package interp

import (
	"fmt"
	"strings"

	"github.com/rmonnet/glox/lang"
)

// Env represents the interpreter state.
// Environment are chained backward to allow lookup
// in enclosing environment (lexical scoping).
type env struct {
	values    map[string]interface{}
	enclosing *env
}

// NewEnv creates a new environment.
func newEnv(enclosing *env) *env {

	return &env{
		values:    make(map[string]interface{}),
		enclosing: enclosing}
}

// define binds a variable name and its value for the environment.
// IfToken the variable was already bound, the name value is bound
// instead and the old value is discarded.
func (e *env) define(name string, value interface{}) {

	e.values[name] = value
}

// get retrieves the value associated with a variable.
// It the variable is not bound a runtimeError is triggered.
func (e *env) get(name *lang.Token) interface{} {

	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	panic(runtimeError{name,
		"Undefined variable '" + name.Lexeme + "'."})
}

// getAt retrieves the value associated with a variable
// in a given enclosing environment. The environment where
// the variable is defined is specified by the distance from
// the current environment.
// There is no error handling because resolver ensure the name
// is in the environment at the proper distance.
func (e *env) getAt(distance int, name string) interface{} {

	return e.ancestor(distance).values[name]
}

// assign binds a new value with an existing variable.
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
		"Undefined variable '" + name.Lexeme + "'."})
}

// assignAt binds a new value with an existing variable,
// looking for the variable in the enclosing environment
// "distance" levels up from the current environment.
// There is no error handling because resolver ensure the name
// is in the environment at the proper distance.
func (e *env) assignAt(distance int, name string, value interface{}) {

	e.ancestor(distance).values[name] = value
}

// ------------------
// Helper Functions
// ------------------

// ancestor return the enclosing environment "distance"
// levels up from the current environment.
func (e *env) ancestor(distance int) *env {

	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

// -----------------
// Debug Functions
// -----------------

// dump print the environment content and enclosing environments
// in the format "distance from current env) key=value".
// It is useful for debugging.
func (e *env) dump(distance int) string {

	b := strings.Builder{}
	for k, v := range e.values {
		fmt.Fprintf(&b, "%d) %s=%v\n", distance, k, v)
	}
	if e.enclosing != nil {
		fmt.Fprint(&b, e.enclosing.dump(distance+1))
	}
	return b.String()
}

// depth returns how many levels down the current environment is
// form the top-level.
// It is useful for debugging.
func (e *env) depth() int {

	i := 0
	for next := e.enclosing; next != nil; next = next.enclosing {
		i++
	}
	return i
}
