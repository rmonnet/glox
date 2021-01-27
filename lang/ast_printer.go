package lang

import "fmt"

// AstPrinter implements a Pretty Printer visitor for lox expression
type AstPrinter struct{}

// Print pretty-prints a lox expression
func (p AstPrinter) Print(e *Expr) {
	(*e).AcceptExpr(p)
}

// VisitLiteral pretty-prints a lox literal expression
func (p AstPrinter) VisitLiteral(l *Literal) {
	fmt.Print(l.Value)
}

// VisitGrouping pretty-prints a lox grouping expression
func (p AstPrinter) VisitGrouping(g *Grouping) {
	fmt.Print("(group ")
	(g.Expression).AcceptExpr(p)
	fmt.Print(")")
}

// VisitUnary pretty-prints a lox unary expression
func (p AstPrinter) VisitUnary(u *Unary) {
	fmt.Print("(")
	fmt.Print(u.Operator.Type)
	(u.Expression).AcceptExpr(p)
	fmt.Print(")")
}

// VisitBinary pretty-prints a lox binary expression
func (p AstPrinter) VisitBinary(b *Binary) {
	fmt.Print("(")
	fmt.Print(b.Operator.Type)
	fmt.Print(" ")
	(b.Left).AcceptExpr(p)
	fmt.Print(" ")
	(b.Right).AcceptExpr(p)
	fmt.Print(")")
}
