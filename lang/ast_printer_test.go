package lang

func ExamplePrettyPrint() {

	e := &BinaryExpr{
		&UnaryExpr{&Token{Minus, "-", 1}, &Lit{123}},
		&Token{Star, "*", 1},
		&GroupingExpr{&Lit{45.67}}}
	PrettyPrint(e)
	// Output: (* (-123) (group 45.67))
}

func ExamplePrettyPrint_string() {

	e := &BinaryExpr{
		&Lit{"abc"},
		&Token{Plus, "+", 1},
		&Lit{"def"}}
	PrettyPrint(e)
	// Output: (+ "abc" "def")
}
