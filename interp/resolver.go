package interp

import (
	"fmt"
	"os"

	"gitlab.com/rcmonnet/glox/lang"
)

// The Resolver type provides operations to resolve variables in
// a lox AST.
type Resolver struct {
	interp               *Interp
	scopes               scopeStack
	currentFunctionScope functionScope
	currentClassScope    classScope
	hadError             bool
}

// NewResolver creates a new resolver and associate it
// with an interpreter.
func NewResolver(i *Interp) *Resolver {

	return &Resolver{interp: i}
}

// resolve goes through an AST tree and resolve variable references.
func (r *Resolver) resolve(statements []lang.Stmt) {

	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

// resolveStmt resolves the variables in the statement.
func (r *Resolver) resolveStmt(stmt lang.Stmt) {

	switch actualStmt := stmt.(type) {
	case *lang.ReturnStmt:
		r.resolveReturnStmt(actualStmt)
	case *lang.PrintStmt:
		r.resolvePrintStmt(actualStmt)
	case *lang.ExprStmt:
		r.resolveExprStmt(actualStmt)
	case *lang.IfStmt:
		r.resolveIfStmt(actualStmt)
	case *lang.WhileStmt:
		r.resolveWhileStmt(actualStmt)
	case *lang.VarDeclStmt:
		r.resolveVarDeclStmt(actualStmt)
	case *lang.ClassDeclStmt:
		r.resolveClassDeclStmt(actualStmt)
	case *lang.FunDeclStmt:
		r.resolveFunDeclStmt(actualStmt)
	case *lang.BlockStmt:
		r.resolveBlockStmt(actualStmt)
	default:
		panic(fmt.Sprintf("Unknown Statement Type: %T", actualStmt))
	}
}

// resolveWhileStmt resolves variables included in a while statement.
func (r *Resolver) resolveWhileStmt(stmt *lang.WhileStmt) {

	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
}

// resolvePrintStmt resolves variables in a print statement.
func (r *Resolver) resolvePrintStmt(stmt *lang.PrintStmt) {

	r.resolveExpr(stmt.Expression)
}

// resolveReturnStmt resolves variables in a return statement.
func (r *Resolver) resolveReturnStmt(stmt *lang.ReturnStmt) {

	// it is an error if returns appears outside of a function
	// definition.
	if r.currentFunctionScope == outsideFunction {
		r.reportError(stmt.Keyword, "Can't return from top-level code.")
	}

	// it is an error to return a value from inside an
	// initializer (initializer always return the class instance).
	if r.currentFunctionScope == inInitializer &&
		stmt.Value != nil {
		r.reportError(stmt.Keyword,
			"Can't return a value from an initializer")
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
}

// resolveExprStmt resolves variables in an expression statement.
func (r *Resolver) resolveExprStmt(stmt *lang.ExprStmt) {

	r.resolveExpr(stmt.Expression)
}

// resolveIfStmt resolves variables in an if statement.
func (r *Resolver) resolveIfStmt(stmt *lang.IfStmt) {

	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}

}

// resolveBlockStmt resolves the variables in the block.
// a block statement represents a new scope/environment.
func (r *Resolver) resolveBlockStmt(stmt *lang.BlockStmt) {

	r.beginScope()
	r.resolve(stmt.Statements)
	r.endScope()
}

// resolveVarDeclStmt resolves a variable declaration.
// This method keeps track of the variable declaration and definition.
func (r *Resolver) resolveVarDeclStmt(stmt *lang.VarDeclStmt) {

	r.declare(stmt.Name)

	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}

	r.define(stmt.Name)
}

// resolveClassDeclStmt resolves a class declaration.
// This method keeps track of the class declaration and definition.
func (r *Resolver) resolveClassDeclStmt(stmt *lang.ClassDeclStmt) {

	enclosingClassScope := r.currentClassScope
	r.currentClassScope = inClass

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	r.scopes.peek()["this"] = true

	for _, method := range stmt.Methods {
		declaration := inMethod
		if method.Name.Lexeme == "init" {
			declaration = inInitializer
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()

	r.currentClassScope = enclosingClassScope
}

// resolveFunDeclStmt resolves a function declaration.
// This method keeps track of the function declaration and definition.
func (r *Resolver) resolveFunDeclStmt(stmt *lang.FunDeclStmt) {

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, inFunction)
}

// resolveFunction resolves variables in a function body.
// The function body represents a new scope/environment.
func (r *Resolver) resolveFunction(stmt *lang.FunDeclStmt, newScope functionScope) {

	enclosingFunctionScope := r.currentFunctionScope
	r.currentFunctionScope = newScope

	r.beginScope()
	for _, param := range stmt.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolve(stmt.Body)
	r.endScope()

	r.currentFunctionScope = enclosingFunctionScope
}

// resolveExpr resolves variable references within an expression.
func (r *Resolver) resolveExpr(expr lang.Expr) {

	switch actualExpr := expr.(type) {
	case *lang.Lit:
		r.resolveLit(actualExpr)
	case *lang.GroupingExpr:
		r.resolveGroupingExpr(actualExpr)
	case *lang.UnaryExpr:
		r.resolveUnaryExpr(actualExpr)
	case *lang.BinaryExpr:
		r.resolveBinaryExpr(actualExpr)
	case *lang.LogicalExpr:
		r.resolveLogicalExpr(actualExpr)
	case *lang.VarExpr:
		r.resolveVarExpr(actualExpr)
	case *lang.AssignExpr:
		r.resolveAssignExpr(actualExpr)
	case *lang.CallExpr:
		r.resolveCallExpr(actualExpr)
	case *lang.ThisExpr:
		r.resolveThisExpr(actualExpr)
	case *lang.GetExpr:
		r.resolveGetExpr(actualExpr)
	case *lang.SetExpr:
		r.resolveSetExpr(actualExpr)
	default:
		panic(fmt.Sprintf("Unknown Expression Type: %T", expr))
	}
}

// resolveUnaryExpr resolves variables in a unary expression.
func (r *Resolver) resolveUnaryExpr(expr *lang.UnaryExpr) {

	r.resolveExpr(expr.Expression)
}

// resolveUnaryExpr resolves variables in a logical expression.
// resolution doesn't consider short-circuit operators (all paths
// must be resolved).
func (r *Resolver) resolveLogicalExpr(expr *lang.LogicalExpr) {

	r.resolveExpr(expr.LeftExpression)
	r.resolveExpr(expr.RightExpression)
}

// resolveLit resolves variables in a literal.
func (r *Resolver) resolveLit(expr *lang.Lit) {
	// nothing to do: there is no variable in a literal.
}

// resolveGroupingExpr resolves variables in a group expression.
func (r *Resolver) resolveGroupingExpr(expr *lang.GroupingExpr) {

	r.resolveExpr(expr.Expression)
}

// resolveGroupingExpr resolves variables in a call expression.
// There is no need to resolve the body of the function at call time.
func (r *Resolver) resolveCallExpr(expr *lang.CallExpr) {

	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}
}

// resolveGetExpr resolves variables in a get expression.
// only the receiver is resolved since fields require dynamic
// dispatch and must be done at runtime.
func (r *Resolver) resolveGetExpr(expr *lang.GetExpr) {

	r.resolveExpr(expr.Object)
}

// resolveSetExpr resolves variables in a set expression.
// only the receiver and the value are resolved since fields
//  require dynamic dispatch and must be done at runtime.
func (r *Resolver) resolveSetExpr(expr *lang.SetExpr) {

	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
}

// resolveBinaryExpr resolves variables in a binary expression.
func (r *Resolver) resolveBinaryExpr(expr *lang.BinaryExpr) {

	r.resolveExpr(expr.LeftExpression)
	r.resolveExpr(expr.RightExpression)
}

// resolveVarExpr resolves variables in a variable expression.
// search for variable definitions in the current scope and
// enclosing scopes.
func (r *Resolver) resolveVarExpr(expr *lang.VarExpr) {

	if !r.scopes.isEmpty() {
		isInitialized, isDefined := r.scopes.peek()[expr.Name.Lexeme]
		if isDefined && !isInitialized {
			r.reportError(expr.Name,
				"Can't read local variable in its own initializer")
		}
	}

	r.resolveLocal(expr, expr.Name)
}

// resolveThisExpr resolves This as a pseudo-variable within
// methods of a class.
func (r *Resolver) resolveThisExpr(expr *lang.ThisExpr) {

	if r.currentClassScope == outsideClass {
		r.reportError(expr.Keyword,
			"can't use 'this' outside of a class.")
	}
	r.resolveLocal(expr, expr.Keyword)
}

// resolveAssignExpr resolves variables in an assignment expression.
// search for variable definitions in the current scope and
// enclosing scopes.
func (r *Resolver) resolveAssignExpr(expr *lang.AssignExpr) {

	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
}

// Helper functions

// beginScope starts a new scope for variable references.
func (r *Resolver) beginScope() {

	r.scopes.push(make(scope))
}

// endScope denotes the end of a scope for variable references.
func (r *Resolver) endScope() {

	r.scopes.pop()
}

// declare associates the variable declaration with the current scope.
// The variable is marked as undefined.
func (r *Resolver) declare(name *lang.Token) {

	if r.scopes.isEmpty() {
		return
	}

	sc := r.scopes.peek()

	// it is an error to redeclare the same variable in the same scope.
	if _, ok := sc[name.Lexeme]; ok {
		r.reportError(name, "Variable already declared in this scope.")
	}

	sc[name.Lexeme] = false
}

// define defines the variable in the current scope.
func (r *Resolver) define(name *lang.Token) {

	if r.scopes.isEmpty() {
		return
	}

	r.scopes.peek()[name.Lexeme] = true
}

// resolveLocal search for the variables in the current scope
// and enclosing scopes and notify the interpreter of the variable
// location.
func (r *Resolver) resolveLocal(expr lang.Expr, name *lang.Token) {

	for i := r.scopes.size() - 1; i >= 0; i-- {
		if _, ok := r.scopes.get(i)[name.Lexeme]; ok {
			r.interp.resolve(expr, r.scopes.size()-1-i)
			return
		}
	}
}

// reportError is triggered when a parser errors is encountered.
// the parser can then continue from that point.
func (r *Resolver) reportError(token *lang.Token, msg string) {

	var where string
	if token.Type == lang.End {
		where = "at end"
	} else {
		where = "at '" + token.Lexeme + "'"
	}
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n",
		token.Line, where, msg)
	r.hadError = true
}

// scope represents an interpreter scope.
type scope map[string]bool

// scopeStack represents a stack of scopes.
type scopeStack struct {
	stack []scope
}

// push pushes a new scope on the stack.
func (s *scopeStack) push(sc scope) {

	s.stack = append(s.stack, sc)
}

// pop returns the latest scope from the stack.
// the latest scope is also removed from the stack.
func (s *scopeStack) pop() scope {

	sc := s.stack[len(s.stack)-1]
	s.stack = s.stack[0 : len(s.stack)-1]
	return sc
}

// peek returns the latest scope from the stack.
// the latest scope is left on the stack.
func (s *scopeStack) peek() scope {

	return s.stack[len(s.stack)-1]
}

// isEmpty checks if the stack is empty.
func (s *scopeStack) isEmpty() bool {

	return len(s.stack) == 0
}

// size returns the number of scopes on the stack.
func (s *scopeStack) size() int {

	return len(s.stack)
}

// get returns the ith scope from the stack.
func (s *scopeStack) get(index int) scope {

	return s.stack[index]
}

// functionScope keeps track if the current scope is a function or
// a method.
type functionScope int

const (
	outsideFunction functionScope = iota
	inFunction
	inInitializer
	inMethod
)

// classScope keeps track if the current scope is within a class.
type classScope int

const (
	outsideClass classScope = iota
	inClass
)
