package lang

func ExamplePrettyPrint() {

	e := &BinaryExpr{
		&UnaryExpr{&Token{MinusToken, "-", 1}, &Lit{123}},
		&Token{StarToken, "*", 1},
		&GroupingExpr{&Lit{45.67}}}
	PrettyPrint(e)
	// Output: (* (-123) (group 45.67))
}

func ExamplePrettyPrint_string() {

	e := &BinaryExpr{
		&Lit{"abc"},
		&Token{PlusToken, "+", 1},
		&Lit{"def"}}
	PrettyPrint(e)
	// Output: (+ "abc" "def")
}
