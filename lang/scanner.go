package lang

// Scanner represents a lox scanner
type Scanner struct {
	source  []rune
	tokens  []*Token
	start   int
	current int
	line    int
}

// NewScanner initialize a new lox scanner.
func NewScanner(source string) *Scanner {

	s := new(Scanner)
	s.source = []rune(source)
	s.line = 1
	// TODO: consider pre-allocating some room in the token array
	return s
}

// ScanTokens scans the source code and return the list
// of tokens.
func (s *Scanner) ScanTokens() []*Token {

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(End, "", "", s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {

	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(RightBrace)
	case '}':
		s.addToken(LeftBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
	case '/':
		if s.match('/') {
			// a comment goes to the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash)
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
			Raise(s.line, "Unexpected character.")
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
		Raise(s.line, "Unterminated string.")
		return
	}

	// need to consume the closing quote
	s.advance()

	// trim the surrounding quotes
	value := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(String, value)
}

// number consumes a number token from the source.
// numbers are integers or simple floating point numbers
// (no exponent). Numbers cannot start or end with a dot.
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

	text := string(s.source[s.start:s.current])
	s.addTokenWithLiteral(Number, text)
}

// identifier consumes an identifier token from the source.
// Identifiers must start with an Alpha character followed
// by any number of AlphaNumeric characters.
func (s *Scanner) identifier() {

	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	tokenType, found := keywords[text]
	if !found {
		tokenType = Identifier
	}
	s.addToken(tokenType)
}

// isDigit checks if the character is a digit
func isDigit(c rune) bool {

	return c >= '0' && c < '9'
}

// isAlpha checks if the character is a letter
// lox only supports ASCII letters
func isAlpha(c rune) bool {

	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') || c == '_'
}

// isAlphaNumeric checks if the character is
// a letter of a digit
func isAlphaNumeric(c rune) bool {

	return isAlpha(c) || isDigit(c)
}

// Helper functions

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
// is as expected. If the character matches, it is consumed.
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
	s.tokens = append(s.tokens, NewToken(tokenType, text, "", s.line))
}

// addTokenWithLiteral adds a token with a literal to the Scanner result
func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal string) {

	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

// keywords is a map including all lox reserved keywords
var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}
