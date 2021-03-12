package interp

import "time"

// lox interpreter built-in functions.
// Each function must implement the loxCallable interface
// (call(), arity()) and the Stringer interface.

// clock represents the built in clock function.
// clock returns the unix time in seconds.
type clock struct{}

// call implements a call to the clock() function.
func (c clock) call(i *Interp, args []interface{}) interface{} {
	return time.Now().Unix()
}

// arity returns the arity of the clock() function.
func (c clock) arity() int {
	return 0
}

// string provides a printable representation of the clock() function.
func (c clock) String() string {
	return "<native fn>"
}
