package lang

// Stmt represents a statement in loc AST.
type Stmt interface {
	stmtNode()
}

// ExprStmt represents an expression statement in loc AST.
type ExprStmt struct {
	Expression Expr
}

// PrintStmt represents a print statement in loc AST.
type PrintStmt struct {
	Expression Expr
}

// BlockStmt represents a block statement in loc AST.
type BlockStmt struct {
	Statements []Stmt
}

// VarDeclStmt represents a variable declaration in loc AST.
type VarDeclStmt struct {
	Name        *Token
	Initializer Expr
}

// Expr represents an expression in loc AST.
type Expr interface {
	exprNode()
}

// Lit represents a STRING, NUMBER, BOOLEAN or NIL literal in loc AST.
type Lit struct {
	Value interface{}
}

// VarExpr represents a variable expression in loc AST.
type VarExpr struct {
	Name *Token
}

// AssignExpr represents an assignment expression in loc AST.
type AssignExpr struct {
	Name  *Token
	Value Expr
}

// BinaryExpr represents a binary expression in loc AST.
type BinaryExpr struct {
	LeftExpression  Expr
	Operator        *Token
	RightExpression Expr
}

// GroupingExpr represents a grouping expression in loc AST.
type GroupingExpr struct {
	Expression Expr
}

// UnaryExpr represents a unary expression in loc AST.
type UnaryExpr struct {
	Operator   *Token
	Expression Expr
}

// Enforce the following types to be Expression.
func (*Lit) exprNode()          {}
func (*UnaryExpr) exprNode()    {}
func (*AssignExpr) exprNode()   {}
func (*BinaryExpr) exprNode()   {}
func (*GroupingExpr) exprNode() {}
func (*VarExpr) exprNode()      {}

// Enforce the following types to be Statement.
func (*ExprStmt) stmtNode()    {}
func (*PrintStmt) stmtNode()   {}
func (*VarDeclStmt) stmtNode() {}
func (*BlockStmt) stmtNode()   {}
