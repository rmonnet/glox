package lang

import "fmt"

// PrettyPrint prints the content of a lox expression as
// a set of nodes.
func PrettyPrint(e Expr) {

	switch n := e.(type) {
	case *StringLit:
		fmt.Printf("%q", n.Value)
	case *NumberLit:
		fmt.Printf("%v", n.Value)
	case *BooleanLit:
		fmt.Printf("%v", n.Value)
	case *NilLit:
		fmt.Print("nil")
	case *GroupingExpr:
		fmt.Print("(group ")
		PrettyPrint(n.Expression)
		fmt.Print(")")
	case *UnaryExpr:
		fmt.Print("(")
		fmt.Print(n.Operator.Type)
		PrettyPrint(n.Expression)
		fmt.Print(")")
	case *BinaryExpr:
		fmt.Print("(")
		fmt.Print(n.Operator.Type)
		fmt.Print(" ")
		PrettyPrint(n.LeftExpression)
		fmt.Print(" ")
		PrettyPrint(n.RightExpression)
		fmt.Print(")")
	default:
		panic(fmt.Sprintf("Unknown Expression Type: %T", e))
	}
}
