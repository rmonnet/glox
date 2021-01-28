package lang

import (
	"fmt"
	"strconv"
	"strings"
)

// ErrParser represents an error occuring in the
// parser.
var ErrParser = fmt.Errorf("parser error")

// Parser represents a lox parser
type Parser struct {
	tokens  []*Token
	current int
}

// NewParser creates a lox parser, using the output
// of a scanner.
func NewParser(tokens []*Token) *Parser {

	p := new(Parser)
	p.tokens = tokens
	return p
}

// Parse parses the stream of tokens into an AST.
func (p *Parser) Parse() (expr Expr) {

	// if an error is reported, need to resynchronize the stream
	defer func() {
		if e := recover(); e != nil {
			if e != ErrParser {
				panic(e)
			}
			// TODO: call synchronize
			expr = nil
		}
	}()

	return p.expression()
}

// Parsing rules

// expression implements the rule for a lox expression.
// expression = equality ;
func (p *Parser) expression() Expr {

	return p.equality()
}

// equality implements the rule for a lox equality expression.
// equality = comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {

	expr := p.comparison()
	for p.match(BangEqual, EqualEqual) {
		op := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{expr, op, right}
	}
	return expr
}

// comparison implements the rule for a lox comparison expression.
// comparison = term ( ( ">" | ">=" | "<" | "<=") term )* ;
func (p *Parser) comparison() Expr {

	expr := p.term()
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		op := p.previous()
		right := p.term()
		expr = &BinaryExpr{expr, op, right}
	}
	return expr
}

// term implements the rule for a lox term expression
// term = factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() Expr {

	expr := p.factor()
	for p.match(Minus, Plus) {
		op := p.previous()
		right := p.factor()
		expr = &BinaryExpr{expr, op, right}
	}
	return expr
}

// factor implements the rule for a lox factor expression
// factor = unary ( ( "/" "*" ) unary )* ;
func (p *Parser) factor() Expr {

	expr := p.unary()
	for p.match(Slash, Star) {
		op := p.previous()
		right := p.unary()
		expr = &BinaryExpr{expr, op, right}
	}
	return expr
}

// unary implements the rule for a lox unary expression
// unary = ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() Expr {

	if p.match(Bang, Equal) {
		op := p.previous()
		right := p.unary()
		return &UnaryExpr{op, right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {

	if p.match(False) {
		return &BooleanLit{false}
	}
	if p.match(True) {
		return &BooleanLit{true}
	}
	if p.match(Nil) {
		return &NilLit{}
	}
	if p.match(Number) {
		n, _ := strconv.ParseFloat(p.previous().Lexeme, 64)
		// TODO: deal with the error in ParseFloat
		// theoretically, there should be no error since
		// we match the token to a float
		return &NumberLit{n}
	}
	if p.match(String) {
		// technically we should remove just a single
		// quote at the beginning and the end of the string
		// but the lox grammar guarantees there is only
		// a single quote at the beginning and end
		s := strings.Trim(p.previous().Lexeme, "\"")
		return &StringLit{s}
	}
	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return &GroupingExpr{expr}
	}
	error(p.peek(), "Expect expression.")
	panic(ErrParser)
}

// Helper functions

// match returns true if the current token matches
// one of the provided token types.
func (p *Parser) match(types ...TokenType) bool {

	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// consume check and skip the next token. If the token
// is different from the expected token, an error is raised.
func (p *Parser) consume(tokenType TokenType, msg string) *Token {

	if p.check(tokenType) {
		return p.advance()
	}

	error(p.peek(), msg)
	panic(ErrParser)
}

// check returns true if the current token matches
// the specified token type.
func (p *Parser) check(tokenType TokenType) bool {

	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

// advance moves to the next token.
func (p *Parser) advance() *Token {

	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if there are more tokens available to parse.
func (p *Parser) isAtEnd() bool {

	return p.peek().Type == End
}

// peek returns the next token in the parsing stream.
func (p *Parser) peek() *Token {

	return p.tokens[p.current]
}

// previous returns the previous token in the parsing stream.
func (p *Parser) previous() *Token {

	return p.tokens[p.current-1]
}

// synchronize search the parsing stream for the first
// token after a semicolon. It is used to continue
// parsing after an error is found and reported.
func synchronize(p *Parser) {

	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}
		switch p.peek().Type {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}
		p.advance()
	}
}

// error is triggered when a parser errors is encountered.
// the parser can then continue from that point.
func error(token *Token, msg string) {

	if token.Type == End {
		report(token.Line, " at end", msg)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", msg)
	}
}
