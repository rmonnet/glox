package lang

// ValueType represents the type of a literal value in lox AST
type ValueType int

const (
	// StringType represents the STRING literal type
	StringType ValueType = iota
	// NumberType represents the NUMBER literal type
	NumberType
	// BooleanType represents the Boolean literal type
	BooleanType
	// NilType represents the Nil literal type
	NilType
)

// Expr represents an expression Node in loc AST
type Expr interface {
	AcceptExpr(v ExprVisitor)
}

// Literal represents a literal Node in loc AST
type Literal struct {
	Type  ValueType
	Value interface{}
}

// AcceptExpr visits a Literal expression
func (l *Literal) AcceptExpr(v ExprVisitor) {
	v.VisitLiteral(l)
}

// Binary represents a binary Node in loc AST
type Binary struct {
	Operator *Token
	Left     Expr
	Right    Expr
}

// AcceptExpr visits a Binary expression
func (b *Binary) AcceptExpr(v ExprVisitor) {
	v.VisitBinary(b)
}

// Grouping represents a grouping Node in loc AST
type Grouping struct {
	Expression Expr
}

// AcceptExpr visits a Grouping expression
func (g *Grouping) AcceptExpr(v ExprVisitor) {
	v.VisitGrouping(g)
}

// Unary represents a unary Node in loc AST
type Unary struct {
	Operator   *Token
	Expression Expr
}

// AcceptExpr visits a unary expression
func (u *Unary) AcceptExpr(v ExprVisitor) {
	v.VisitUnary(u)
}

// ExprVisitor represents a visitor object on the lox Expressions
type ExprVisitor interface {
	VisitLiteral(node *Literal)
	VisitUnary(node *Unary)
	VisitBinary(node *Binary)
	VisitGrouping(node *Grouping)
}
