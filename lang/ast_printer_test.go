package lang

func ExamplePrettyPrint() {

	e := &BinaryExpr{
		&Token{Star, "*", "", 1},
		&UnaryExpr{&Token{Minus, "-", "", 1}, &NumberLit{123}},
		&GroupingExpr{&NumberLit{45.67}}}
	PrettyPrint(e)
	// Output: (* (-123) (group 45.67))
}
