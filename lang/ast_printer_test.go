package lang

func ExampleExpr() {

	p := AstPrinter{}
	e := Binary{
		&Token{Star, "*", "", 1},
		&Unary{&Token{Minus, "-", "", 1}, &Literal{NumberType, 123}},
		&Grouping{&Literal{NumberType, 45.67}}}
	e.AcceptExpr(p)
	// Output: (* (-123) (group 45.67))
}
