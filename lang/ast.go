package lang

// Stmt represents a statement in lox AST.
type Stmt interface {
	stmtNode()
}

// ExprStmt represents an expression statement in lox AST.
type ExprStmt struct {
	Expression Expr
}

// ClassDeclStmt represents a class definition in lox AST.
type ClassDeclStmt struct {
	Name    *Token
	Methods []*FunDeclStmt
}

// FunDeclStmt represents a function definition in lox AST.
type FunDeclStmt struct {
	Name   *Token
	Params []*Token
	Body   []Stmt
}

// IfStmt represents an if statement in lox AST.
type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

// PrintStmt represents a print statement in lox AST.
type PrintStmt struct {
	Expression Expr
}

// BlockStmt represents a block statement in lox AST.
type BlockStmt struct {
	Statements []Stmt
}

// VarDeclStmt represents a variable declaration in lox AST.
type VarDeclStmt struct {
	Name        *Token
	Initializer Expr
}

// WhileStmt represents a while statement in lox AST.
type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

// ReturnStmt represents a return statement in lox AST.
type ReturnStmt struct {
	Keyword *Token
	Value   Expr
}

// Expr represents an expression in lox AST.
type Expr interface {
	exprNode()
}

// Lit represents a STRING, NUMBER, BOOLEAN or NIL literal in lox AST.
type Lit struct {
	Value interface{}
}

// VarExpr represents a variable expression in lox AST.
type VarExpr struct {
	Name *Token
}

// AssignExpr represents an assignment expression in lox AST.
type AssignExpr struct {
	Name  *Token
	Value Expr
}

// BinaryExpr represents a binary expression in lox AST.
type BinaryExpr struct {
	LeftExpression  Expr
	Operator        *Token
	RightExpression Expr
}

// LogicalExpr represents a logical expression in lox AST.
type LogicalExpr struct {
	LeftExpression  Expr
	Operator        *Token
	RightExpression Expr
}

// SetExpr represents read write to a class field in lox AST.
type SetExpr struct {
	Object Expr
	Name   *Token
	Value  Expr
}

// ThisExpr represents the pseudo-variable this representing
// a class instance in lox AST.
type ThisExpr struct {
	Keyword *Token
}

// CallExpr represents a function call in lox AST.
type CallExpr struct {
	Callee    Expr
	Paren     *Token
	Arguments []Expr
}

// GetExpr represents read access to a class field in lox AST.
type GetExpr struct {
	Object Expr
	Name   *Token
}

// GroupingExpr represents a grouping expression in lox AST.
type GroupingExpr struct {
	Expression Expr
}

// UnaryExpr represents a unary expression in lox AST.
type UnaryExpr struct {
	Operator   *Token
	Expression Expr
}

// Enforce the following types to be Expression.
func (*Lit) exprNode()          {}
func (*UnaryExpr) exprNode()    {}
func (*AssignExpr) exprNode()   {}
func (*BinaryExpr) exprNode()   {}
func (*LogicalExpr) exprNode()  {}
func (*CallExpr) exprNode()     {}
func (*GetExpr) exprNode()      {}
func (*SetExpr) exprNode()      {}
func (*ThisExpr) exprNode()     {}
func (*GroupingExpr) exprNode() {}
func (*VarExpr) exprNode()      {}

// Enforce the following types to be Statement.
func (*ExprStmt) stmtNode()      {}
func (*ClassDeclStmt) stmtNode() {}
func (*FunDeclStmt) stmtNode()   {}
func (*IfStmt) stmtNode()        {}
func (*PrintStmt) stmtNode()     {}
func (*ReturnStmt) stmtNode()    {}
func (*WhileStmt) stmtNode()     {}
func (*VarDeclStmt) stmtNode()   {}
func (*BlockStmt) stmtNode()     {}
