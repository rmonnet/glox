package lang

import (
	"strings"
	"testing"
)

func TestScanTokens(t *testing.T) {

	script :=
		`and ! != class , . else	= == false fun for > >=	an_Identifier01
	if { ( < <= - nil 123 123.456 or + print return } ) ; / *
	"a string" super this true var while
	// a comment`

	expect := []string{
		"and", "!", "!=", "class", ",", ".", "else", "=", "==",
		"false", "fun", "for", ">", ">=",
		"Identifier(an_Identifier01)", "if", "{", "(", "<", "<=",
		"-", "nil", "Number(123)", "Number(123.456)", "or", "+",
		"print", "return", "}", ")", ";", "/", "*", "String(a string)",
		"super", "this", "true", "var", "while", "end-of-stream"}

	matchTokens(t, expect, script)
}

func TestScanNumbers(t *testing.T) {

	t.Run("Parse integer", func(t *testing.T) {

		scanValidToken(t, "Number(1234)", "1234")
	})

	t.Run("Parse float", func(t *testing.T) {

		scanValidToken(t, "Number(12.349)", "12.349")
	})

	// Note: floats starting with '.' or ending with '.'
	// like '.1234' nd '1234.' in lox are scanned as
	// 2 tokens (a number and a dot).
	// The fact that they are invalid is reported by the parser
	// (i.e they are interpreted as a dot call with an integer
	// receiver or target).

	t.Run("Parse float starting with a dot", func(t *testing.T) {

		expect := []string{".", "Number(1234)", "end-of-stream"}
		matchTokens(t, expect, ".1234")
	})

	t.Run("Parse float ending with a dot", func(t *testing.T) {

		expect := []string{"Number(1234)", ".", "end-of-stream"}
		matchTokens(t, expect, "1234.")
	})

}

func TestScanStrings(t *testing.T) {

	t.Run("Parse regular string", func(t *testing.T) {

		scanValidToken(t, "String(hello world)", "\"hello world\"")
	})

	t.Run("Parse multiline string", func(t *testing.T) {

		scanValidToken(t, "String(hello\nworld)", "\"hello\nworld\"")
	})

	t.Run("Parse unterminated string", func(t *testing.T) {

		scanInvalidToken(t, "\"helloworld")
	})

}

// ------------------
// Helper functions
// ------------------

func matchTokens(t *testing.T, expect []string, script string) {

	t.Helper()

	scanner := &Scanner{}
	got := scanner.ScanTokens(script)

	if scanner.HadError() {
		t.Error("Error encountered while scanning")
	}

	length := len(expect)
	if len(got) > length {
		length = len(got)
	}

	for i := 0; i < length; i++ {

		if i >= len(got) {
			t.Errorf("Expected token '%s' was missing in %dth position",
				expect[i], i+1)
		} else if i >= len(expect) {
			t.Errorf("Unexpected token '%s' in %dth position",
				got[i], i+1)
		} else if got[i].String() != expect[i] {
			t.Errorf("Expected token '%s' but got '%s' in %dth position",
				expect[i], got[i], i+1)
		}
	}
}

func scanValidToken(t *testing.T, expect string, script string) {

	scanner := &Scanner{}
	tokens := scanner.ScanTokens(script)
	if scanner.HadError() {
		t.Error("Error encountered while scanning")
	} else if len(tokens) != 2 {
		// there is always an extra end-of-stream token
		t.Errorf("Expected 1 token but got %d", len(tokens))
	} else if tokens[0].String() != expect {
		t.Errorf("Expected token '%s' but got '%s", expect, tokens[0])
	}

}

func scanInvalidToken(t *testing.T, script string) {

	scanner := &Scanner{}
	errOut := &strings.Builder{}
	scanner.RedirectErrors(errOut)
	tokens := scanner.ScanTokens(script)
	if !scanner.HadError() {
		t.Error("Expect Error was not reported by scanner")
	} else if len(tokens) != 1 {
		// there is always an extra end-of-stream token
		t.Errorf("Expected 0 token but got %d", len(tokens))
	}

}
