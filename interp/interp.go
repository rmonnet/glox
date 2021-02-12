// Package interp implements the tree-walker
// interpreter for the lox language
package interp

import (
	"fmt"
	"reflect"

	"gitlab.com/rcmonnet/glox/lang"
)

// Interp represents the state of the lox interpreter.
type Interp struct {
	hadCompileError bool
	hadRuntimeError bool
	env             *env
}

// New creates a new interpreter.
func New() *Interp {

	interp := &Interp{}
	interp.env = newEnv(nil)
	return interp
}

// Run runs the lox interpreter on the provided
// program.
func (i *Interp) Run(script string) {

	scanner := lang.NewScanner(script)
	tokens := scanner.ScanTokens()

	parser := lang.NewParser(tokens)
	statements := parser.Parse()

	if scanner.HadError() || parser.HadError() {
		i.hadCompileError = true
		return
	}

	i.interpret(statements)
}

// HadCompileError indicates if errors occurred during
// compilation.
func (i *Interp) HadCompileError() bool {

	return i.hadCompileError
}

// HadRuntimeError indicates if errors occurred during
// compilation.
func (i *Interp) HadRuntimeError() bool {

	return i.hadRuntimeError
}

// runtimeError represents an error encountered during
// Runtime interpretation.
type runtimeError struct {
	token   *lang.Token
	message string
}

// Error extracts the Error Message out of a runtimeError.
func (e runtimeError) Error() string {
	return e.message
}

// interpret evaluates the expression and display the result.
func (i *Interp) interpret(statements []lang.Stmt) {
	defer func() {
		if e := recover(); e != nil {
			rte := e.(runtimeError)
			fmt.Printf("%s\n[line %d]\n", rte.message, rte.token.Line)
			i.hadRuntimeError = true
		}
	}()
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

// execute executes a statement.
func (i *Interp) execute(stmt lang.Stmt) {

	switch s := stmt.(type) {
	case *lang.PrintStmt:
		i.executePrintStmt(s)
	case *lang.ExprStmt:
		i.executeExprStmt(s)
	case *lang.IfStmt:
		i.executeIfStmt(s)
	case *lang.WhileStmt:
		i.executeWhileStmt(s)
	case *lang.VarDeclStmt:
		i.executeValDeclStmt(s)
	case *lang.BlockStmt:
		i.executeBlockStmt(s)
	default:
		panic(fmt.Sprintf("Unknown Statement Type: %T", s))
	}
}

// executeWhileStmt executes a while statement.
func (i *Interp) executeWhileStmt(stmt *lang.WhileStmt) {

	for isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
}

// executeIfStmt executes an if statement.
func (i *Interp) executeIfStmt(stmt *lang.IfStmt) {

	if isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

// executeBlockStmt executes a block statement.
func (i *Interp) executeBlockStmt(block *lang.BlockStmt) {

	previousEnv := i.env
	// ensure that the previous environment is restored
	// no matter what happens.
	defer func() {
		i.env = previousEnv
	}()
	i.env = newEnv(previousEnv)
	for _, s := range block.Statements {
		i.execute(s)
	}
}

// executeExprstmt executes an expression statement.
func (i *Interp) executeExprStmt(stmt *lang.ExprStmt) {

	i.evaluate(stmt.Expression)
}

// executePrintStmt executes a print statement.
func (i *Interp) executePrintStmt(stmt *lang.PrintStmt) {

	val := i.evaluate(stmt.Expression)
	fmt.Println(stringify(val))
}

func (i *Interp) executeValDeclStmt(stmt *lang.VarDeclStmt) {

	var val interface{}
	if stmt.Initializer != nil {
		val = i.evaluate(stmt.Initializer)
	}

	i.env.define(stmt.Name.Lexeme, val)
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
func (i *Interp) evaluate(e lang.Expr) interface{} {

	switch n := e.(type) {
	case *lang.Lit:
		return n.Value
	case *lang.GroupingExpr:
		return i.evaluate(n.Expression)
	case *lang.UnaryExpr:
		return i.evaluateUnary(n)
	case *lang.BinaryExpr:
		return i.evaluateBinary(n)
	case *lang.LogicalExpr:
		return i.evaluateLogical(n)
	case *lang.VarExpr:
		return i.env.get(n.Name)
	case *lang.AssignExpr:
		return i.evaluateAssign(n)
	default:
		panic(fmt.Sprintf("Unknown Expression Type: %T", e))
	}
}

// evaluateLogical evaluates a Logical expression and return
// the result as a literal
func (i *Interp) evaluateLogical(l *lang.LogicalExpr) interface{} {

	left := i.evaluate(l.LeftExpression)
	if l.Operator.Type == lang.Or {
		if isTruthy(left) {
			return left
		}
	} else if l.Operator.Type == lang.And {
		if !isTruthy(left) {
			return left
		}
	} else {
		panic(fmt.Sprintf("Unknown Logical Operator %v",
			l.Operator))
	}
	return i.evaluate(l.RightExpression)
}

// evaluateAssign evaluates an Assignment expression and returns
// the result as a literal
func (i *Interp) evaluateAssign(a *lang.AssignExpr) interface{} {

	value := i.evaluate(a.Value)
	i.env.assign(a.Name, value)
	return value
}

// evaluateUnary evaluates a Unary expression and returns
// the result as a literal
func (i *Interp) evaluateUnary(u *lang.UnaryExpr) interface{} {

	right := i.evaluate(u.Expression)
	switch u.Operator.Type {
	case lang.Minus:
		val := operandToNumber(u.Operator, right)
		return -val
	case lang.Bang:
		return !isTruthy(right)
	}
	return nil
}

func (i *Interp) evaluateBinary(b *lang.BinaryExpr) interface{} {

	left := i.evaluate(b.LeftExpression)
	right := i.evaluate(b.RightExpression)

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
		// to make it easier to debug,
		// when used for string concatenation, "+" supports
		// implicit conversion to string
		if leftVal.Kind() == reflect.String || rightVal.Kind() == reflect.String {
			return valueToString(leftVal) + valueToString(rightVal)
		}
		panic(runtimeError{b.Operator,
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
		panic(runtimeError{operator, "Operand must be a number."})
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
		panic(runtimeError{operator, "Operands must be numbers."})
	}
	return leftVal, rightVal
}

// valueToString converts any of the lox primitive types
// to a string. It is used for implicit conversion to
// string for the "+" operator.
func valueToString(value reflect.Value) string {

	kind := value.Kind()
	switch {
	case kind == reflect.String:
		return value.String()
	case kind == reflect.Float64:
		return fmt.Sprintf("%v", value.Float())
	case kind == reflect.Bool:
		return fmt.Sprintf("%v", value.Bool())
	case !value.IsValid():
		return "nil"
	default:
		panic(fmt.Sprintf("Unexpected primitive type %s, %v", kind, value))
	}
}
