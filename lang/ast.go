package lang

import (
	"fmt"
	"strings"
)

// ------------
// Statements
// ------------

type PrettyPrinter interface {
	PrettyPrint(pad, tab string) string
}

// Stmt represents a statement in lox AST.
type Stmt interface {
	PrettyPrinter
	fmt.Stringer
	stmtNode()
}

// BlockStmt represents a block statement in lox AST.
type BlockStmt struct {
	Statements []Stmt
}

func (*BlockStmt) stmtNode() {}

func (stmt *BlockStmt) PrettyPrint(pad, tab string) string {

	b := strings.Builder{}
	fmt.Fprintf(&b, "%s(block", pad)
	newPad := pad + tab
	for _, stmt := range stmt.Statements {
		fmt.Fprintf(&b, "%s", stmt.PrettyPrint(newPad, tab))
	}
	fmt.Fprint(&b, ")")
	return b.String()
}

func (stmt *BlockStmt) String() string {

	b := strings.Builder{}
	fmt.Fprintf(&b, "(block")
	for _, stmt := range stmt.Statements {
		fmt.Fprintf(&b, " %s", stmt.String())
	}
	fmt.Fprint(&b, ")")
	return b.String()
}

// ClassDeclStmt represents a class definition in lox AST.
type ClassDeclStmt struct {
	Name       *Token
	Superclass *VarExpr
	Methods    []*FunDeclStmt
}

func (*ClassDeclStmt) stmtNode() {}

func (stmt *ClassDeclStmt) PrettyPrint(pad, tab string) string {

	b := strings.Builder{}
	if stmt.Superclass != nil {
		fmt.Fprintf(&b, "%s(class %s %s", pad, stmt.Name.Lexeme,
			stmt.Superclass.Name.Lexeme)
	} else {
		fmt.Fprintf(&b, "%s(class %s nil", pad, stmt.Name.Lexeme)
	}
	newPad := pad + tab
	for _, method := range stmt.Methods {
		fmt.Fprintf(&b, "%s", method.PrettyPrint(newPad, tab))
	}
	fmt.Fprint(&b, ")")
	return b.String()
}

func (stmt *ClassDeclStmt) String() string {

	b := strings.Builder{}
	if stmt.Superclass != nil {
		fmt.Fprintf(&b, "(class %s %s", stmt.Name.Lexeme,
			stmt.Superclass.Name.Lexeme)
	} else {
		fmt.Fprintf(&b, "(class %s nil", stmt.Name.Lexeme)
	}
	for _, method := range stmt.Methods {
		fmt.Fprintf(&b, " %s", method.String())
	}
	fmt.Fprint(&b, ")")
	return b.String()
}

// ExprStmt represents an expression statement in lox AST.
type ExprStmt struct {
	Expression Expr
}

func (*ExprStmt) stmtNode() {}

func (stmt *ExprStmt) PrettyPrint(pad, tab string) string {

	return fmt.Sprintf("%s%s", pad, stmt.Expression.String())
}

func (stmt *ExprStmt) String() string {

	return stmt.Expression.String()

}

// FunDeclStmt represents a function definition in lox AST.
type FunDeclStmt struct {
	Name   *Token
	Params []*Token
	Body   []Stmt
}

func (*FunDeclStmt) stmtNode() {}

func (stmt *FunDeclStmt) PrettyPrint(pad, tab string) string {

	b := strings.Builder{}
	fmt.Fprintf(&b, "%s(fun %s (params", pad, stmt.Name.Lexeme)
	for _, param := range stmt.Params {
		fmt.Fprintf(&b, " %s", param.Lexeme)
	}
	fmt.Fprint(&b, ")")
	newPad := pad + tab
	for _, statement := range stmt.Body {
		fmt.Fprintf(&b, "%s", statement.PrettyPrint(newPad, tab))
	}
	fmt.Fprint(&b, ")")

	return b.String()
}

func (stmt *FunDeclStmt) String() string {

	b := strings.Builder{}
	fmt.Fprintf(&b, "(fun %s (params", stmt.Name.Lexeme)
	for _, param := range stmt.Params {
		fmt.Fprintf(&b, " %s", param.Lexeme)
	}
	fmt.Fprint(&b, ")")
	for _, statement := range stmt.Body {
		fmt.Fprintf(&b, " %s", statement.String())
	}
	fmt.Fprint(&b, ")")
	return b.String()
}

// IfStmt represents an if statement in lox AST.
type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (*IfStmt) stmtNode() {}

func (stmt *IfStmt) PrettyPrint(pad, tab string) string {

	b := strings.Builder{}
	newPad := pad + tab
	fmt.Fprintf(&b, "%s(if %s%s", pad, stmt.Condition.String(),
		stmt.ThenBranch.PrettyPrint(newPad, tab))
	if stmt.ElseBranch != nil {
		fmt.Fprintf(&b, "%s", stmt.ElseBranch.PrettyPrint(newPad, tab))
	}
	fmt.Fprint(&b, ")")

	return b.String()
}

func (stmt *IfStmt) String() string {

	b := strings.Builder{}
	fmt.Fprintf(&b, "(if %s %s", stmt.Condition.String(),
		stmt.ThenBranch.String())
	if stmt.ElseBranch != nil {
		fmt.Fprintf(&b, " %s", stmt.ElseBranch.String())
	}
	fmt.Fprint(&b, ")")
	return b.String()
}

// PrintStmt represents a print statement in lox AST.
type PrintStmt struct {
	Expression Expr
}

func (*PrintStmt) stmtNode() {}

func (stmt *PrintStmt) PrettyPrint(pad, tab string) string {

	return fmt.Sprintf("%s(print %s)", pad, stmt.Expression.String())
}

func (stmt *PrintStmt) String() string {

	return fmt.Sprintf("(print %s)", stmt.Expression.String())
}

// ReturnStmt represents a return statement in lox AST.
type ReturnStmt struct {
	Keyword *Token
	Value   Expr
}

func (*ReturnStmt) stmtNode() {}

func (stmt *ReturnStmt) PrettyPrint(pad, tab string) string {

	if stmt.Value != nil {
		return fmt.Sprintf("%s(return %s)", pad, stmt.Value.String())
	} else {
		return fmt.Sprintf("%s(return)", pad)
	}
}

func (stmt *ReturnStmt) String() string {

	if stmt.Value != nil {
		return fmt.Sprintf("(return %s)", stmt.Value.String())
	} else {
		return fmt.Sprintf("(return)")
	}
}

// VarDeclStmt represents a variable declaration in lox AST.
type VarDeclStmt struct {
	Name        *Token
	Initializer Expr
}

func (*VarDeclStmt) stmtNode() {}

func (stmt *VarDeclStmt) PrettyPrint(pad, tab string) string {

	if stmt.Initializer != nil {
		return fmt.Sprintf("%s(var %s %s)", pad, stmt.Name.Lexeme,
			stmt.Initializer.String())
	} else {
		return fmt.Sprintf("%s(var %s)", pad, stmt.Name.Lexeme)
	}
}

func (stmt *VarDeclStmt) String() string {

	if stmt.Initializer != nil {
		return fmt.Sprintf("(var %s %s)", stmt.Name.Lexeme,
			stmt.Initializer.String())
	} else {
		return fmt.Sprintf("(var %s)", stmt.Name.Lexeme)
	}
}

// WhileStmt represents a while statement in lox AST.
type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (*WhileStmt) stmtNode() {}

func (stmt *WhileStmt) PrettyPrint(pad, tab string) string {

	return fmt.Sprintf("%s(while %s%s)", pad,
		stmt.Condition.String(), stmt.Body.PrettyPrint(pad+tab, tab))
}

func (stmt *WhileStmt) String() string {

	return fmt.Sprintf("(while %s %s)",
		stmt.Condition.String(), stmt.Body.String())
}

// -------------
// Expressions
// -------------

// Expr represents an expression in lox AST.
type Expr interface {
	fmt.Stringer
	exprNode()
}

// AssignExpr represents an assignment expression in lox AST.
type AssignExpr struct {
	Name  *Token
	Value Expr
}

func (*AssignExpr) exprNode() {}

func (expr *AssignExpr) String() string {

	return fmt.Sprintf("(assign %s %s)", expr.Name.Lexeme,
		expr.Value)
}

// BinaryExpr represents a binary expression in lox AST.
type BinaryExpr struct {
	LeftExpression  Expr
	Operator        *Token
	RightExpression Expr
}

func (*BinaryExpr) exprNode() {}

func (expr *BinaryExpr) String() string {

	return fmt.Sprintf("(%s %s %s)", expr.Operator.Lexeme,
		expr.LeftExpression.String(), expr.RightExpression.String())
}

// CallExpr represents a function call in lox AST.
type CallExpr struct {
	Callee    Expr
	Paren     *Token
	Arguments []Expr
}

func (*CallExpr) exprNode() {}

func (expr *CallExpr) String() string {

	b := strings.Builder{}
	fmt.Fprintf(&b, "(call %s (args", expr.Callee.String())
	for _, argument := range expr.Arguments {
		fmt.Fprintf(&b, " %s", argument.String())
	}
	fmt.Fprint(&b, "))")
	return b.String()
}

// GetExpr represents read access to a class field in lox AST.
type GetExpr struct {
	Object Expr
	Name   *Token
}

func (*GetExpr) exprNode() {}

func (expr *GetExpr) String() string {

	return fmt.Sprintf("(get %s %s)", expr.Object.String(),
		expr.Name.Lexeme)
}

// GroupingExpr represents a grouping expression in lox AST.
type GroupingExpr struct {
	Expression Expr
}

func (*GroupingExpr) exprNode() {}

func (expr *GroupingExpr) String() string {

	return fmt.Sprintf("(group %s)", expr.Expression)
}

// Lit represents a STRING, NUMBER, BOOLEAN or NIL literal in lox AST.
type Lit struct {
	Value interface{}
}

func (*Lit) exprNode() {}

func (expr *Lit) String() string {

	if expr.Value == nil {
		return "nil"
	}
	if s, ok := expr.Value.(string); ok {
		return fmt.Sprintf("\"%s\"", s)
	}
	return fmt.Sprintf("%v", expr.Value)
}

// LogicalExpr represents a logical expression in lox AST.
type LogicalExpr struct {
	LeftExpression  Expr
	Operator        *Token
	RightExpression Expr
}

func (*LogicalExpr) exprNode() {}

func (expr *LogicalExpr) String() string {

	return fmt.Sprintf("(%s %s %s)", expr.Operator.Lexeme,
		expr.LeftExpression.String(), expr.RightExpression.String())
}

// SetExpr represents read write to a class field in lox AST.
type SetExpr struct {
	Object Expr
	Name   *Token
	Value  Expr
}

func (*SetExpr) exprNode() {}

func (expr *SetExpr) String() string {

	return fmt.Sprintf("(set %s %s %s)", expr.Object.String(),
		expr.Name.Lexeme, expr.Value.String())
}

// SuperExpr represents the pseudo-variable "super" representing
// a class superclass in lox AST.
type SuperExpr struct {
	Keyword *Token
	Method  *Token
}

func (*SuperExpr) exprNode() {}

func (expr *SuperExpr) String() string {

	return fmt.Sprintf("(super %s)", expr.Method.Lexeme)
}

// ThisExpr represents the pseudo-variable "this" representing
// a class instance in lox AST.
type ThisExpr struct {
	Keyword *Token
}

func (*ThisExpr) exprNode() {}

func (expr *ThisExpr) String() string {

	return "(this)"
}

// UnaryExpr represents a unary expression in lox AST.
type UnaryExpr struct {
	Operator   *Token
	Expression Expr
}

func (*UnaryExpr) exprNode() {}

func (expr *UnaryExpr) String() string {

	return fmt.Sprintf("(%s %s)", expr.Operator.Lexeme,
		expr.Expression.String())
}

// VarExpr represents a variable expression in lox AST.
type VarExpr struct {
	Name *Token
}

func (*VarExpr) exprNode() {}

func (expr *VarExpr) String() string {

	return fmt.Sprintf("(%s)", expr.Name.Lexeme)
}
