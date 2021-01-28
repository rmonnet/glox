package lang

func ExamplePrettyPrint() {

	e := &BinaryExpr{
		&UnaryExpr{&Token{Minus, "-", 1}, &NumberLit{123}},
		&Token{Star, "*", 1},
		&GroupingExpr{&NumberLit{45.67}}}
	PrettyPrint(e)
	// Output: (* (-123) (group 45.67))
}

func ExamplePrettyPrint_string() {

	e := &BinaryExpr{
		&StringLit{"abc"},
		&Token{Star, "+", 1},
		&StringLit{"def"}}
	PrettyPrint(e)
	// Output: (+ "abc" "def)
}
