// Package interp implements the tree-walker
// interpreter for the lox language
package interp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"gitlab.com/rcmonnet/glox/lang"
)

const (
	exDataErr = 65
	exSwErr   = 70
)

// hadRuntimeError indicates an error occurred while
// interpreting the lox script.
var hadRuntimeError bool

// RuntimeError represents an error encountered during
// RUntime interpretation.
type RuntimeError struct {
	token   *lang.Token
	message string
}

// Error extracts the Error Message out of a RuntimeError.
func (e RuntimeError) Error() string {
	return e.message
}

// RunFile runs the lox interpreter on the
// script in the file
func RunFile(filename string) {

	script, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("unable to read ", filename)
		os.Exit(exDataErr)
	}
	run(string(script))
	if lang.HadError {
		os.Exit(exDataErr)
	}
	if hadRuntimeError {
		os.Exit(exSwErr)
	}
}

// RunPrompt runs the lox interpreter interactively
func RunPrompt() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Println("")
			break
		}
		run(scanner.Text())
		lang.HadError = false
		hadRuntimeError = false
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error while reading ", err)
		os.Exit(exDataErr)
	}

}

// run runs the lox interpreter on the provided
// program.
func run(script string) {

	scanner := lang.NewScanner(script)
	tokens := scanner.ScanTokens()

	// for debugging only
	for _, token := range tokens {
		fmt.Println(token)
	}

	parser := lang.NewParser(tokens)
	expr := parser.Parse()

	if lang.HadError {
		return
	}

	// for debugging only
	lang.PrettyPrint(expr)
	fmt.Println("")

	// TODO: we will need an interpreter object
	// to store state
	interpret(expr)
}

// interpret evaluates the expression and display the result.
func interpret(e lang.Expr) {
	defer func() {
		if e := recover(); e != nil {
			rte := e.(RuntimeError)
			fmt.Printf("%s\n[line %d]\n", rte.message, rte.token.Line)
			hadRuntimeError = true
		}
	}()
	val := evaluate(e)
	fmt.Println(stringify(val))
}

// stringify returns a valid lox string representation
// of the literal.
func stringify(lit interface{}) string {

	if lit == nil {
		return "nil"
	}
	// original code remove ".0" suffix from floats
	// to show they represent integers. Go '%v'
	// does this automatically
	return fmt.Sprintf("%v", lit)
}

// evaluate evaluates an expression and returns the result
// as a literal
func evaluate(e lang.Expr) interface{} {

	switch n := e.(type) {
	case *lang.Lit:
		return n.Value
	case *lang.GroupingExpr:
		return evaluate(n.Expression)
	case *lang.UnaryExpr:
		return evaluateUnary(n)
	case *lang.BinaryExpr:
		return evaluateBinary(n)
	default:
		panic(fmt.Sprintf("Unknown Expression Type: %T", e))
	}
}

// evaluateUnary evaluates a Unary expression and returns
// the result as a literal
func evaluateUnary(u *lang.UnaryExpr) interface{} {

	right := evaluate(u.Expression)
	switch u.Operator.Type {
	case lang.Minus:
		val := operandToNumber(u.Operator, right)
		return -val
	case lang.Bang:
		return !isTruthy(right)
	}
	return nil
}

func evaluateBinary(b *lang.BinaryExpr) interface{} {

	left := evaluate(b.LeftExpression)
	right := evaluate(b.RightExpression)

	switch b.Operator.Type {
	case lang.Minus:
		leftVal, rightVal := operandsToNumbers(
			b.Operator, left, right)
		return leftVal - rightVal
	case lang.Slash:
		leftVal, rightVal := operandsToNumbers(
			b.Operator, left, right)
		return leftVal / rightVal
	case lang.Star:
		leftVal, rightVal := operandsToNumbers(
			b.Operator, left, right)
		return leftVal * rightVal
	case lang.Plus:
		leftVal := reflect.ValueOf(left)
		rightVal := reflect.ValueOf(right)
		if leftVal.Kind() == reflect.Float64 && rightVal.Kind() == reflect.Float64 {
			return leftVal.Float() + rightVal.Float()
		}
		if leftVal.Kind() == reflect.String && rightVal.Kind() == reflect.String {
			return leftVal.String() + rightVal.String()
		}
		panic(RuntimeError{b.Operator,
			"Operands must be two numbers or two strings."})
	case lang.Greater:
		leftVal, rightVal := operandsToNumbers(
			b.Operator, left, right)
		return leftVal > rightVal
	case lang.GreaterEqual:
		leftVal, rightVal := operandsToNumbers(
			b.Operator, left, right)
		return leftVal >= rightVal
	case lang.Less:
		leftVal, rightVal := operandsToNumbers(
			b.Operator, left, right)
		return leftVal < rightVal
	case lang.LessEqual:
		leftVal, rightVal := operandsToNumbers(
			b.Operator, left, right)
		return leftVal <= rightVal
	case lang.BangEqual:
		return !isEqual(left, right)
	case lang.EqualEqual:
		return isEqual(left, right)
	}
	return nil
}

// isTruthy evaluate if the literal is true.
// In lox, false and nil are false, everything else is true
func isTruthy(lit interface{}) bool {
	if lit == nil {
		return false
	}
	if val, ok := lit.(bool); ok {
		return val
	}
	return true
}

// isEqual checks if two lox literals are equal
func isEqual(left interface{}, right interface{}) bool {

	// comparing incomparable types in go may cause a panic
	// but at this point left and right can only be
	// lox literals, that is NUMBER, STRING or BOOLEAN
	return left == right
}

// operandToNumber convert the operand to a lox number
// or panic if the type is incorrect.
func operandToNumber(operator *lang.Token,
	operand interface{}) float64 {

	val, ok := operand.(float64)
	if !ok {
		panic(RuntimeError{operator, "Operand must be a number."})
	}
	return val
}

// operandsToNumbers converts the operands to lox numbers
// or panic if the types are incorrect.
func operandsToNumbers(operator *lang.Token,
	left, right interface{}) (float64, float64) {

	leftVal, leftOk := left.(float64)
	rightVal, rightOk := right.(float64)
	if !leftOk || !rightOk {
		panic(RuntimeError{operator, "Operands must be numbers."})
	}
	return leftVal, rightVal
}
