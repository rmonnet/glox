package lang

// Expr represents an expression in loc AST
type Expr interface {
	exprNode()
}

// Lit represents a STRING, NUMBER, BOOLEAN or NIL literal in loc AST
type Lit struct {
	Value interface{}
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
func (*Lit) exprNode()          {}
func (*UnaryExpr) exprNode()    {}
func (*BinaryExpr) exprNode()   {}
func (*GroupingExpr) exprNode() {}
