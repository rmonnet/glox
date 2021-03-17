package lang

import (
	"fmt"
	"io"
	"os"
)

// Scanner represents a lox scanner.
type Scanner struct {
	source   []rune
	tokens   []*Token
	start    int
	current  int
	line     int
	hadError bool
	errOut   io.Writer
}

// RedirectErrors switches the file errors are written to.
// Errors go to stderr by default.
func (s *Scanner) RedirectErrors(errOut io.Writer) {

	s.errOut = errOut
}

// ScanTokens scans the source code and return the list
// of tokens.
func (s *Scanner) ScanTokens(source string) []*Token {

	// Reset the scanner state in case it is reused.
	s.source = []rune(source)
	s.tokens = nil
	s.start = 0
	s.current = 0
	s.line = 1
	s.hadError = false
	if s.errOut == nil {
		s.errOut = os.Stderr
	}

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, &Token{EndToken, "", s.line})
	return s.tokens
}

// HadError reports if some errors were encountered during
// scanning. It should be called after ScanTokens before using
// the result.
func (s *Scanner) HadError() bool {

	return s.hadError
}

// scanToken scans the new token in the script.
func (s *Scanner) scanToken() {

	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParenToken)
	case ')':
		s.addToken(RightParenToken)
	case '{':
		s.addToken(LeftBraceToken)
	case '}':
		s.addToken(RightBraceToken)
	case ',':
		s.addToken(CommaToken)
	case '.':
		s.addToken(DotToken)
	case '-':
		s.addToken(MinusToken)
	case '+':
		s.addToken(PlusToken)
	case ';':
		s.addToken(SemicolonToken)
	case '*':
		s.addToken(StarToken)
	case '!':
		if s.match('=') {
			s.addToken(BangEqualToken)
		} else {
			s.addToken(BangToken)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqualToken)
		} else {
			s.addToken(EqualToken)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqualToken)
		} else {
			s.addToken(LessToken)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqualToken)
		} else {
			s.addToken(GreaterToken)
		}
	case '/':
		if s.match('/') {
			// a comment goes to the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SlashToken)
		}
	case ' ', '\r', '\t':
		// ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.reportError("Unexpected character.")
			// TODO: it would be nicer to coalesce all the consecutive erroneous characters
			// into a single error message
		}
	}
}

// string consumes a string token from the source.
// strings are defined using double quotes.
// lox supports multilines strings.
func (s *Scanner) string() {

	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.reportError("Unterminated string.")
		return
	}

	// need to consume the closing quote
	s.advance()

	s.addToken(StringToken)
}

// number consumes a number token from the source.
// numbers are integers or simple floating point numbers
// (no exponent). Numbers cannot start or end with a dot,
// in that case, they will be parsed as two tokens (a number)
// and a dot).
func (s *Scanner) number() {

	for isDigit(s.peek()) {
		s.advance()
	}

	// look for the fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
	}

	for isDigit(s.peek()) {
		s.advance()
	}

	s.addToken(NumberToken)
}

// identifier consumes an identifier token from the source.
// Identifiers must start with an Alpha character followed
// by any number of AlphaNumeric characters.
func (s *Scanner) identifier() {

	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IdentifierToken
	}

	s.addToken(tokenType)
}

// isDigit checks if the character is a digit.
func isDigit(c rune) bool {

	return c >= '0' && c <= '9'
}

// isAlpha checks if the character is a letter.
// Lox only supports ASCII letters.
func isAlpha(c rune) bool {

	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') || c == '_'
}

// isAlphaNumeric checks if the character is
// a letter of a digit
func isAlphaNumeric(c rune) bool {

	return isAlpha(c) || isDigit(c)
}

// ------------------
// Helper functions
// ------------------

// reportError reports an error during interpretation
func (s *Scanner) reportError(message string) {

	fmt.Fprintf(s.errOut, "[line %d] Error: %s\n",
		s.line, message)
	s.hadError = true
}

// isAtEnd checks if the scanner has reached the end of the
// source file.
func (s *Scanner) isAtEnd() bool {

	return s.current >= len(s.source)
}

// advance advances by one character in the source
func (s *Scanner) advance() rune {

	s.current++
	return s.source[s.current-1]
}

// match checks the next character in the source
// is as expected. IfToken the character matches, it is consumed.
func (s *Scanner) match(expected rune) bool {

	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

// peek returns the next character in the source but
// doesn't advance the counter
func (s *Scanner) peek() rune {

	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current]
}

// peekNext returns the second character ahead in the
// source but doesn't advance the counter
func (s *Scanner) peekNext() rune {

	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

// addToken adds a token to the Scanner result
func (s *Scanner) addToken(tokenType TokenType) {

	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, &Token{tokenType, text, s.line})
}

// keywords is a map including all lox reserved keywords
var keywords = map[string]TokenType{
	"and":    AndToken,
	"class":  ClassToken,
	"else":   ElseToken,
	"false":  FalseToken,
	"for":    ForToken,
	"fun":    FunToken,
	"if":     IfToken,
	"nil":    NilToken,
	"or":     OrToken,
	"print":  PrintToken,
	"return": ReturnToken,
	"super":  SuperToken,
	"this":   ThisToken,
	"true":   TrueToken,
	"var":    VarToken,
	"while":  WhileToken,
}
