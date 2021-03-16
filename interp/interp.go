// Package interp implements the tree-walker
// interpreter for the lox language
package interp

import (
	"fmt"
	"io"
	"os"

	"gitlab.com/rcmonnet/glox/lang"
)

// Interp represents the state of the lox interpreter.
type Interp struct {
	hadCompileError bool
	hadRuntimeError bool
	globalEnv       *env
	env             *env
	locals          map[lang.Expr]int
	out             io.Writer
	errOut          io.Writer
}

// New creates a new interpreter.
func New(out, errOut io.Writer) *Interp {

	interp := &Interp{}
	interp.globalEnv = newEnv(nil)
	interp.globalEnv.define("clock", clock{})
	interp.env = interp.globalEnv
	interp.locals = make(map[lang.Expr]int)
	if out == nil {
		interp.out = os.Stdout
	} else {
		interp.out = out
	}
	if errOut == nil {
		interp.errOut = os.Stderr
	} else {
		interp.errOut = errOut
	}
	return interp
}

// Run runs the lox interpreter on the provided program.
func (i *Interp) Run(script string) {

	scanner := &lang.Scanner{}
	scanner.RedirectErrors(i.errOut)
	tokens := scanner.ScanTokens(script)

	parser := &lang.Parser{}
	parser.RedirectErrors(i.errOut)
	statements := parser.Parse(tokens)

	if scanner.HadError() || parser.HadError() {
		i.hadCompileError = true
		return
	}

	resolver := NewResolver(i)
	resolver.RedirectErrors(i.errOut)
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
// ThisToken is used in conjunction with panic to unwind the stack
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
	case *lang.ClassDeclStmt:
		i.executeClassDeclStmt(actualStmt)
	case *lang.FunDeclStmt:
		i.executeFunDeclStmt(actualStmt)
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
	fmt.Fprintln(i.out, stringify(value))
}

// executeValDeclStmt executes a variable declaration.
func (i *Interp) executeValDeclStmt(stmt *lang.VarDeclStmt) {

	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.env.define(stmt.Name.Lexeme, value)
}

// executeClassDeclStmt executes a class declaration.
func (i *Interp) executeClassDeclStmt(stmt *lang.ClassDeclStmt) {

	// the variable referencing the superclass must evaluate to
	// a class.
	var superclass *loxClass
	if stmt.Superclass != nil {
		sc := i.evaluate(stmt.Superclass)
		var ok bool
		if superclass, ok = sc.(*loxClass); !ok {
			panic(runtimeError{stmt.Superclass.Name,
				"Superclass must be a class."})
		}
	}

	// separate definition from assignment to allow
	// reference to the class inside its own methods.
	i.env.define(stmt.Name.Lexeme, nil)

	environment := i.env
	if stmt.Superclass != nil {
		environment = newEnv(i.env)
		environment.define("super", superclass)
	}

	methods := make(map[string]*loxFunction)
	for _, method := range stmt.Methods {
		isInitializer := method.Name.Lexeme == "init"
		function := &loxFunction{method, environment, isInitializer}
		methods[method.Name.Lexeme] = function
	}

	class := &loxClass{stmt.Name.Lexeme, superclass, methods}

	i.env.assign(stmt.Name, class)
}

// executeFunDeclStmt executes a function declaration.
func (i *Interp) executeFunDeclStmt(stmt *lang.FunDeclStmt) {

	function := &loxFunction{stmt, i.env, false}
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
		return i.evaluateVar(actualExpr)
	case *lang.ThisExpr:
		return i.evaluateThis(actualExpr)
	case *lang.SuperExpr:
		return i.evaluateSuper(actualExpr)
	case *lang.AssignExpr:
		return i.evaluateAssign(actualExpr)
	case *lang.CallExpr:
		return i.evaluateCall(actualExpr)
	case *lang.GetExpr:
		return i.evaluateGet(actualExpr)
	case *lang.SetExpr:
		return i.evaluateSet(actualExpr)
	default:
		panic(fmt.Sprintf("Unknown Expression Type: %T", expr))
	}
}

// evaluateVar evaluates a variable and returns its value.
func (i *Interp) evaluateVar(expr *lang.VarExpr) interface{} {

	return i.lookupVariable(expr.Name, expr)
}

// evaluateThis evaluates the "this" pseudo-variable and returns
// the instance it is pointing to.
func (i *Interp) evaluateThis(expr *lang.ThisExpr) interface{} {

	return i.lookupVariable(expr.Keyword, expr)
}

// evaluateSuper evaluates the "super" pseudo-variable and returns
// the method in the super class it is pointing to.
func (i *Interp) evaluateSuper(expr *lang.SuperExpr) interface{} {

	distance := i.locals[expr]
	superclass := i.env.getAt(distance, "super").(*loxClass)

	// we need to bound the method to 'this' in the 'calling' environment
	// not in the 'super' environment.
	// 'this' environment is always directly below 'super' environment.
	this := i.env.getAt(distance-1, "this").(*loxInstance)

	method, ok := superclass.findMethod(expr.Method.Lexeme)
	if ok {
		return method.bind(this)
	}

	panic(runtimeError{expr.Method,
		fmt.Sprintf("Undefined method '%s'.", expr.Method.Lexeme)})

}

// evaluateLogical evaluates a Logical expression and return
// the result as a literal.
// Logical operators implements short-circuits (if the result
// can be determined from the left operand, the right one is not
// evaluated).
func (i *Interp) evaluateLogical(expr *lang.LogicalExpr) interface{} {

	left := i.evaluate(expr.LeftExpression)

	switch expr.Operator.Type {
	case lang.OrToken:
		if isTruthy(left) {
			return left
		}
	case lang.AndToken:
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
	case lang.MinusToken:
		val := toNumber(expr.Operator, right)
		return -val
	case lang.BangToken:
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
	case lang.MinusToken:
		return toNumber(op, left) - toNumber(op, right)
	case lang.SlashToken:
		return toNumber(op, left) / toNumber(op, right)
	case lang.StarToken:
		return toNumber(op, left) * toNumber(op, right)
	case lang.PlusToken:
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
	case lang.GreaterToken:
		return toNumber(op, left) > toNumber(op, right)
	case lang.GreaterEqualToken:
		return toNumber(op, left) >= toNumber(op, right)
	case lang.LessToken:
		return toNumber(op, left) < toNumber(op, right)
	case lang.LessEqualToken:
		return toNumber(op, left) <= toNumber(op, right)
	case lang.BangEqualToken:
		return !isEqual(left, right)
	case lang.EqualEqualToken:
		return isEqual(left, right)
	}
	return nil
}

// evaluateCall evaluates a function calls and return the
// result as a literal.
func (i *Interp) evaluateCall(c *lang.CallExpr) interface{} {

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

// evaluateGet evaluates a field reference and return the
// result as a literal.
func (i *Interp) evaluateGet(expr *lang.GetExpr) interface{} {

	object := i.evaluate(expr.Object)

	instance, ok := object.(*loxInstance)

	if !ok {
		panic(runtimeError{expr.Name,
			"Only class instances have fields."})
	}

	return instance.get(expr.Name)
}

// evaluateSet assigns a field reference and return the
// assigned value as a literal.
func (i *Interp) evaluateSet(expr *lang.SetExpr) interface{} {

	object := i.evaluate(expr.Object)

	instance, ok := object.(*loxInstance)

	if !ok {
		panic(runtimeError{expr.Name,
			"Only class instances have fields."})
	}

	value := i.evaluate(expr.Value)

	instance.set(expr.Name, value)
	return value
}

// --------------------------------
// functions and class structures
// --------------------------------

// the loxCallable interface represents a lox function or closure.
type loxCallable interface {
	call(*Interp, []interface{}) interface{}
	arity() int
}

// the loxFunction represents non-native lox functions.
type loxFunction struct {
	decl          *lang.FunDeclStmt
	closure       *env
	isInitializer bool
}

// call evaluates the body of a lox function.
func (f *loxFunction) call(interp *Interp, args []interface{}) (result interface{}) {

	// intercept panic returning a returnValue.
	// this is used by the return statement to ensure
	// the stack is properly unwound regardless of how
	// deeply nested the return statement is.
	defer func() {
		if err := recover(); err != nil {
			if retval, ok := err.(returnValue); ok {
				// initializer always return class instance.
				if f.isInitializer {
					result = f.closure.getAt(0, "this")
				} else {
					result = retval.value
				}
			} else {
				panic(err)
			}
		}
	}()

	env := newEnv(f.closure)

	for i := 0; i < len(f.decl.Params); i++ {
		env.define(f.decl.Params[i].Lexeme, args[i])
	}

	interp.executeBlockStmt(f.decl.Body, env)

	// "init()" always returns a reference to the class instance,
	// even if called directly.
	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}
	return nil
}

// arity returns the number of parameters expected by a lox function.
func (f *loxFunction) arity() int {

	return len(f.decl.Params)
}

// bind creates a new function with the same body but
// a new environment with a bound value of "this".
// It ties a method to the specific class instance
// it references.
func (f *loxFunction) bind(instance *loxInstance) *loxFunction {

	env := newEnv(f.closure)
	env.define("this", instance)
	return &loxFunction{f.decl, env, f.isInitializer}
}

// string returns a string representation of a lox function.
func (f *loxFunction) String() string {

	return fmt.Sprintf("<fun %s>", f.decl.Name.Lexeme)
}

type loxClass struct {
	Name       string
	Superclass *loxClass
	Methods    map[string]*loxFunction
}

// call creates an instance of a lox class.
func (c *loxClass) call(interp *Interp, args []interface{}) interface{} {

	instance := newLoxInstance(c)

	if initializer, ok := c.findMethod("init"); ok {
		initializer.bind(instance).call(interp, args)
	}

	return instance
}

// arity returns the number of parameters expected by a lox class
// constructor.
func (c *loxClass) arity() int {

	if initializer, ok := c.findMethod("init"); ok {
		return initializer.arity()
	}

	return 0
}

// findMethod look up the requested method name in the class.
func (c *loxClass) findMethod(name string) (*loxFunction, bool) {

	method, ok := c.Methods[name]
	if ok {
		return method, true
	}

	if c.Superclass != nil {
		return c.Superclass.findMethod(name)
	}

	return nil, false
}

// string returns a string representation of a lox class.
func (c *loxClass) String() string {

	return fmt.Sprintf("<class %s>", c.Name)
}

// loxInstance represents an instance of a lox class.
type loxInstance struct {
	class  *loxClass
	fields map[string]interface{}
}

// newLoxInstance creates a new instance of the given class.
func newLoxInstance(class *loxClass) *loxInstance {

	instance := &loxInstance{
		class:  class,
		fields: make(map[string]interface{}),
	}
	return instance
}

// get retrieves the value associated with the instance field
// or raise an error if the field is undefined.
func (i *loxInstance) get(name *lang.Token) interface{} {

	// lookup name can be a field or a method
	value, ok := i.fields[name.Lexeme]

	if ok {
		return value
	}

	method, ok := i.class.findMethod(name.Lexeme)

	if ok {
		return method.bind(i)
	}

	panic(runtimeError{name,
		fmt.Sprintf("Undefined field or method '%s'.", name.Lexeme)})
}

// set assigns a value to an instance field. IfToken this field
// is undefined, set adds it to the instance.
func (i *loxInstance) set(name *lang.Token, value interface{}) {

	i.fields[name.Lexeme] = value
}

// string returns a string representation of a lox instance.
func (i *loxInstance) String() string {

	return fmt.Sprintf("<instance %s>", i.class.Name)
}

// ------------------
// Helper functions
// ------------------

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

	// TODO: it should be sufficient to just printf("%v", value)
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
	case *loxFunction:
		return v.String()
	case *loxClass:
		return v.String()
	case *loxInstance:
		return v.String()
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
