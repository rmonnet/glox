// Package interp implements the tree-walker
// interpreter for the lox language
package interp

import (
	"fmt"

	"gitlab.com/rcmonnet/glox/lang"
)

// the loxCallable interface represents a lox function or closure.
type loxCallable interface {
	call(*Interp, []interface{}) interface{}
	arity() int
}

// the loxFunction represents non-native lox functions.
type loxFunction struct {
	decl    *lang.FunStmt
	closure *env
}

// call evaluates the body of a lox function.
func (f *loxFunction) call(i *Interp, args []interface{}) interface{} {

	env := newEnv(f.closure)

	for i := 0; i < len(f.decl.Params); i++ {
		env.define(f.decl.Params[i].Lexeme, args[i])
	}

	i.executeBlockStmt(f.decl.Body, env)
	return nil
}

// arity returns the number of parameters expected by a lox function.
func (f *loxFunction) arity() int {

	return len(f.decl.Params)
}

// string returns a string representation of a lox function.
func (f *loxFunction) String() string {

	return fmt.Sprintf("<fn %s>", f.decl.Name.Lexeme)
}

// Interp represents the state of the lox interpreter.
type Interp struct {
	hadCompileError bool
	hadRuntimeError bool
	globalEnv       *env
	env             *env
	locals          map[lang.Expr]int
}

// New creates a new interpreter.
func New() *Interp {

	interp := &Interp{}
	interp.globalEnv = newEnv(nil)
	interp.globalEnv.define("clock", clock{})
	interp.env = interp.globalEnv
	interp.locals = make(map[lang.Expr]int)
	return interp
}

// Run runs the lox interpreter on the provided program.
func (i *Interp) Run(script string) {

	scanner := lang.NewScanner(script)
	tokens := scanner.ScanTokens()

	parser := lang.NewParser(tokens)
	statements := parser.Parse()

	if scanner.HadError() || parser.HadError() {
		i.hadCompileError = true
		return
	}

	resolver := NewResolver(i)
	resolver.resolve(statements)

	if resolver.hadError {
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

// returnValue represents a return object.
// This is used in conjunction with panic to unwind the stack
// to the point of the function call and return the value.
type returnValue struct {
	value interface{}
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

	switch actualStmt := stmt.(type) {
	case *lang.ReturnStmt:
		i.executeReturnStmt(actualStmt)
	case *lang.PrintStmt:
		i.executePrintStmt(actualStmt)
	case *lang.ExprStmt:
		i.executeExprStmt(actualStmt)
	case *lang.IfStmt:
		i.executeIfStmt(actualStmt)
	case *lang.WhileStmt:
		i.executeWhileStmt(actualStmt)
	case *lang.VarDeclStmt:
		i.executeValDeclStmt(actualStmt)
	case *lang.FunStmt:
		i.executeFunStmt(actualStmt)
	case *lang.BlockStmt:
		i.executeBlockStmt(actualStmt.Statements, newEnv(i.env))
	default:
		panic(fmt.Sprintf("Unknown Statement Type: %T", stmt))
	}
}

// executeWhileStmt executes a while statement.
func (i *Interp) executeWhileStmt(stmt *lang.WhileStmt) {

	for isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
}

func (i *Interp) executeReturnStmt(stmt *lang.ReturnStmt) {

	var value interface{}
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	// here panic is used in an exception-like pattern
	// to unwind the stack up to the call return point.
	panic(returnValue{value})
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
// We are passing the set of statement directly so we
// can reuse that method to execute a function body during a call.
func (i *Interp) executeBlockStmt(statements []lang.Stmt, blockEnv *env) {

	previousEnv := i.env

	// ensure that the previous environment is restored
	// no matter what happens.
	defer func() {
		i.env = previousEnv
	}()

	i.env = blockEnv
	for _, s := range statements {
		i.execute(s)
	}
}

// executeExprstmt executes an expression statement.
func (i *Interp) executeExprStmt(stmt *lang.ExprStmt) {

	i.evaluate(stmt.Expression)
}

// executePrintStmt executes a print statement.
func (i *Interp) executePrintStmt(stmt *lang.PrintStmt) {

	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))
}

// executeValDeclStmt executes a variable declaration.
func (i *Interp) executeValDeclStmt(stmt *lang.VarDeclStmt) {

	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.env.define(stmt.Name.Lexeme, value)
}

// executeFunDeclStmt executes a function declaration.
func (i *Interp) executeFunStmt(stmt *lang.FunStmt) {

	function := &loxFunction{stmt, i.env}
	i.env.define(stmt.Name.Lexeme, function)
}

// evaluate evaluates an expression and returns the result
// as a literal
func (i *Interp) evaluate(expr lang.Expr) interface{} {

	switch actualExpr := expr.(type) {
	case *lang.Lit:
		return actualExpr.Value
	case *lang.GroupingExpr:
		return i.evaluate(actualExpr.Expression)
	case *lang.UnaryExpr:
		return i.evaluateUnary(actualExpr)
	case *lang.BinaryExpr:
		return i.evaluateBinary(actualExpr)
	case *lang.LogicalExpr:
		return i.evaluateLogical(actualExpr)
	case *lang.VarExpr:
		return i.evaluateVarExpr(actualExpr)
	case *lang.AssignExpr:
		return i.evaluateAssign(actualExpr)
	case *lang.CallExpr:
		return i.evaluateCall(actualExpr)
	default:
		panic(fmt.Sprintf("Unknown Expression Type: %T", expr))
	}
}

// evaluateVarExpr evaluate a variable and returns its value.
func (i *Interp) evaluateVarExpr(expr *lang.VarExpr) interface{} {

	return i.lookupVariable(expr.Name, expr)
}

// evaluateLogical evaluates a Logical expression and return
// the result as a literal.
// Logical operators implements short-circuits (if the result
// can be determined from the left operand, the right one is not
// evaluated).
func (i *Interp) evaluateLogical(expr *lang.LogicalExpr) interface{} {

	left := i.evaluate(expr.LeftExpression)

	switch expr.Operator.Type {
	case lang.Or:
		if isTruthy(left) {
			return left
		}
	case lang.And:
		if !isTruthy(left) {
			return left
		}
	default:
		panic(fmt.Sprintf("Unknown Logical Operator %v",
			expr.Operator))
	}
	return i.evaluate(expr.RightExpression)
}

// evaluateAssign evaluates an Assignment expression and returns
// the result as a literal.
func (i *Interp) evaluateAssign(expr *lang.AssignExpr) interface{} {

	value := i.evaluate(expr.Value)
	i.assignVariable(expr, value)
	return value
}

// evaluateUnary evaluates a Unary expression and returns
// the result as a literal.
func (i *Interp) evaluateUnary(expr *lang.UnaryExpr) interface{} {

	right := i.evaluate(expr.Expression)

	switch expr.Operator.Type {
	case lang.Minus:
		val := toNumber(expr.Operator, right)
		return -val
	case lang.Bang:
		return !isTruthy(right)
	default:
		return nil
	}
}

// evaluateBinary evaluates a Binary expresion and returns the
// result as a literal.
func (i *Interp) evaluateBinary(expr *lang.BinaryExpr) interface{} {

	left := i.evaluate(expr.LeftExpression)
	right := i.evaluate(expr.RightExpression)
	op := expr.Operator

	switch op.Type {
	case lang.Minus:
		return toNumber(op, left) - toNumber(op, right)
	case lang.Slash:
		return toNumber(op, left) / toNumber(op, right)
	case lang.Star:
		return toNumber(op, left) * toNumber(op, right)
	case lang.Plus:
		if isNumber(left) && isNumber(right) {
			return toNumber(op, left) + toNumber(op, right)
		}
		// to make it easier to debug,
		// when used for string concatenation, "+" supports
		// implicit conversion to string
		if isString(left) || isString(right) {
			return toString(left) + toString(right)
		}
		panic(runtimeError{expr.Operator,
			"Operands must be two numbers or at least one string."})
	case lang.Greater:
		return toNumber(op, left) > toNumber(op, right)
	case lang.GreaterEqual:
		return toNumber(op, left) >= toNumber(op, right)
	case lang.Less:
		return toNumber(op, left) < toNumber(op, right)
	case lang.LessEqual:
		return toNumber(op, left) <= toNumber(op, right)
	case lang.BangEqual:
		return !isEqual(left, right)
	case lang.EqualEqual:
		return isEqual(left, right)
	}
	return nil
}

// evaluateCall evaluates a function calls and return the
// result as a literal
func (i *Interp) evaluateCall(c *lang.CallExpr) (result interface{}) {

	// intercept panic returning a returnValue.
	// this is used by the return statement to ensure
	// the stack is properly unwound regardless of how
	// deeply nested the return statement is.
	defer func() {
		if err := recover(); err != nil {
			if retval, ok := err.(returnValue); ok {
				result = retval.value
			} else {
				panic(err)
			}
		}
	}()

	callee := i.evaluate(c.Callee)

	var arguments []interface{}
	for _, arg := range c.Arguments {
		arguments = append(arguments, i.evaluate(arg))
	}

	function, ok := callee.(loxCallable)
	if !ok {
		panic(runtimeError{c.Paren, "Can only call functions and classes."})
	}
	if len(arguments) != function.arity() {
		panic(runtimeError{c.Paren, fmt.Sprintf(
			"Expected %d arguments but got %d.", function.arity(), len(arguments))})
	}
	return function.call(i, arguments)
}

// Helper functions

// resolve keep track of which environment the expression
// is defined in.
// It is called by the Resolver static analyzer.
func (i *Interp) resolve(expr lang.Expr, depth int) {

	i.locals[expr] = depth
}

// lookupVariable looks up the specific variable in the
// environment using lexical scoping.
// The specific environment level to select was specified
// by the static analyzer using the resolve method.
func (i *Interp) lookupVariable(name *lang.Token, expr lang.Expr) interface{} {

	if distance, ok := i.locals[expr]; ok {
		return i.env.getAt(distance, name.Lexeme)
	}
	return i.globalEnv.get(name)
}

// assignVariable assign the specified value to the variable
// in the environment using lexical scoping.
// The specific environment level to select was specified
// by the static analyzer using the resolve method.
func (i *Interp) assignVariable(expr *lang.AssignExpr, value interface{}) {

	if distance, ok := i.locals[expr]; ok {
		i.env.assignAt(distance, expr.Name, value)
	} else {
		i.globalEnv.assign(expr.Name, value)
	}
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

// toNumber convert the operand to a lox number
// or panic if the type is incorrect.
func toNumber(operator *lang.Token,
	operand interface{}) float64 {

	val, ok := operand.(float64)
	if !ok {
		panic(runtimeError{operator, "Operand must be a number."})
	}
	return val
}

// toString converts any of the lox primitive types
// to a string. It is used for implicit conversion to
// string for the "+" operator.
func toString(value interface{}) string {

	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%v", v)
	default:
		panic(fmt.Sprintf("Unexpected primitive type %T", value))
	}
}

// isNumber checks if a generic interface represents a lox float.
func isNumber(value interface{}) bool {

	_, ok := value.(float64)
	return ok
}

// isString checks if a generic interface represents a lox float.
func isString(value interface{}) bool {

	_, ok := value.(string)
	return ok
}
