package lang

// Expr represents an expression in loc AST
type Expr interface {
	exprNode()
}

// StringLit represents a STRING literal in loc AST
type StringLit struct {
	Value string
}

// NumberLit represents a NUMBER literal in loc AST
type NumberLit struct {
	Value float64
}

// BooleanLit represents a BOOLEAN literal in loc AST
type BooleanLit struct {
	Value bool
}

// NilLit represents the NIL literal in loc AST
type NilLit struct {
}

// BinaryExpr represents a binary expression in loc AST
type BinaryExpr struct {
	LeftExpression  Expr
	Operator        *Token
	RightExpression Expr
}

// GroupingExpr represents a grouping expression in loc AST
type GroupingExpr struct {
	Expression Expr
}

// UnaryExpr represents a unary expression in loc AST
type UnaryExpr struct {
	Operator   *Token
	Expression Expr
}

// Enforce the following types to be Expression
func (*StringLit) exprNode()    {}
func (*NumberLit) exprNode()    {}
func (*BooleanLit) exprNode()   {}
func (*NilLit) exprNode()       {}
func (*UnaryExpr) exprNode()    {}
func (*BinaryExpr) exprNode()   {}
func (*GroupingExpr) exprNode() {}
