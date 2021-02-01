package lang

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// errParser is a marker for an error occuring in the
// parser. It is used to trigger synchronization.
var errParser = fmt.Errorf("parser error")

// Parser represents a lox parser
type Parser struct {
	tokens   []*Token
	current  int
	hadError bool
}

// NewParser creates a lox parser, using the output
// of a scanner.
func NewParser(tokens []*Token) *Parser {

	p := new(Parser)
	p.tokens = tokens
	return p
}

// Parse parses the stream of tokens into an AST.
func (p *Parser) Parse() []Stmt {

	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements

}

// HadError reports if some errors were encountered during
// the parsing phase. It should be checked before the
// result is used.
func (p *Parser) HadError() bool {

	return p.hadError
}

// Parsing rules

// declaration implements the rule for a lox declaration.
// declaration = varDeclStmt | statement ;
func (p *Parser) declaration() (statement Stmt) {

	// if an error is reported while parsing a declaration
	// or a statement, need to resynchronize the stream
	defer func() {
		if e := recover(); e != nil {
			if e != errParser {
				panic(e)
			}
			p.synchronize()
			statement = nil
		}
	}()

	if p.match(Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

// varDeclaration implements the rule for a lox variable declaration.
// varDeclStmt = "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDeclaration() Stmt {

	name := p.consume(Identifier, "Expect variable name.")
	var initializer Expr
	if p.match(Equal) {
		initializer = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after variable declaration.")
	return &VarDeclStmt{name, initializer}

}

// statement implements the rule for a lox statement.
// statement = printStmt | exprStmt | block ;
func (p *Parser) statement() Stmt {

	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(LeftBrace) {
		return p.blockStatement()
	}
	return p.expressionStatement()

}

// blockStatement implements the rule for a lox block.
// block = "{" declaration* "}" ;
func (p *Parser) blockStatement() Stmt {

	var statements []Stmt
	for !p.check(RightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(RightBrace, "Expect '}' after block.")
	return &BlockStmt{statements}
}

// printStatement implements the rule for a lox PrintStmt.
// printStmt = "print" expression ";" ;
func (p *Parser) printStatement() Stmt {

	expr := p.expression()
	p.consume(Semicolon, "Expect ';' after value.")
	return &PrintStmt{expr}
}

// expressionStatement implements the rule for a lox exprStmt
// exprStmt = expression ";" ;
func (p *Parser) expressionStatement() Stmt {

	expr := p.expression()
	p.consume(Semicolon, "Expect ';' after expression.")
	return &ExprStmt{expr}
}

// expression implements the rule for a lox expression.
// expression = assignment ;
func (p *Parser) expression() Expr {

	return p.assignment()
}

// assignment implements the rule for a lox assignment expression.
// assignment = IDENTIFIER "=" assignment | equality ;
func (p *Parser) assignment() Expr {

	// Because we may need an infinite look-ahead to find the "=" token
	// we treat the left side as any expression and only
	// check if it is an identifier when we find the "=" token.

	expr := p.equality()
	if p.match(Equal) {
		equals := p.previous()
		value := p.assignment()
		if varExpr, ok := expr.(*VarExpr); ok {
			return &AssignExpr{varExpr.Name, value}
		}
		p.reportError(equals, "Invalid assignment target.")
	}
	return expr
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

	if p.match(Bang, Minus) {
		op := p.previous()
		right := p.unary()
		return &UnaryExpr{op, right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {

	if p.match(False) {
		return &Lit{false}
	}
	if p.match(True) {
		return &Lit{true}
	}
	if p.match(Nil) {
		return &Lit{}
	}
	if p.match(Number) {
		n, _ := strconv.ParseFloat(p.previous().Lexeme, 64)
		// TODO: deal with the error in ParseFloat
		// theoretically, there should be no error since
		// we match the token to a float
		return &Lit{n}
	}
	if p.match(String) {
		// technically we should remove just a single
		// quote at the beginning and the end of the string
		// but the lox grammar guarantees there is only
		// a single quote at the beginning and end
		s := strings.Trim(p.previous().Lexeme, "\"")
		return &Lit{s}
	}
	if p.match(Identifier) {
		return &VarExpr{p.previous()}
	}
	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return &GroupingExpr{expr}
	}
	p.reportError(p.peek(), "Expect expression.")
	panic(errParser)
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

	p.reportError(p.peek(), msg)
	panic(errParser)
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
func (p *Parser) synchronize() {

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

// reportError is triggered when a parser errors is encountered.
// the parser can then continue from that point.
func (p *Parser) reportError(token *Token, msg string) {

	var where string
	if token.Type == End {
		where = "at end"
	} else {
		where = "at '" + token.Lexeme + "'"
	}
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n",
		token.Line, where, msg)
	p.hadError = true
}