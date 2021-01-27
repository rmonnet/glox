// Package lang contains definition for the lox
// language grammar.
package lang

import "fmt"

// TokenType represents the type of a lox token.
type TokenType int

const (
	// LeftParen represents a '(' token.
	LeftParen TokenType = iota
	// RightParen represents a ')' token.
	RightParen
	// LeftBrace represents a '{' token.
	LeftBrace
	// RightBrace represents a '}' token.
	RightBrace
	// Comma represents a ',' token.
	Comma
	// Dot represents a '.' token.
	Dot
	// Minus represents a '-' token.
	Minus
	// Plus represents a '+' token.
	Plus
	// Semicolon represents a ';' token.
	Semicolon
	// Slash represents a '/' token.
	Slash
	// Star represents a '*' token.
	Star
	// Bang represents a '!' token.
	Bang
	// BangEqual represents a '!=' token.
	BangEqual
	// Equal represents an '=' token.
	Equal
	// EqualEqual represents an '==' token.
	EqualEqual
	// Greater represents a '>' token.
	Greater
	// GreaterEqual represents a '>=' token.
	GreaterEqual
	// Less represents a '<'' token.
	Less
	// LessEqual represents a '<=' token.
	LessEqual
	// Identifier represents any identifier token.
	Identifier
	// String represents any string token.
	String
	// Number represents any number token.
	Number
	// And represents an 'and' token.
	And
	// Class represents a 'class' token.
	Class
	// Else represents an 'else' token.
	Else
	// False represents a 'false' token.
	False
	// Fun represents a 'fun' token.
	Fun
	// For represents a 'for' token.
	For
	// If represents an 'if' token.
	If
	// Nil represents a 'nil' token.
	Nil
	// Or represents an 'or' token.
	Or
	// Print represents a 'print' token.
	Print
	// Return represents a 'return' token.
	Return
	// Super represents a 'super' token.
	Super
	// This represents a 'this' token.
	This
	// True represents a 'true' token.
	True
	// Var represents a 'var' token.
	Var
	// While represents a 'while' token.
	While
	// End is a special token that represents the end of stream.
	End
)

// Token represents a lox token.
type Token struct {
	Type   TokenType
	Lexeme string
	// TODO: do we need to store the literal?
	Literal string
	Line    int
}

// NewToken creates a new token.
func NewToken(tokenType TokenType, lexeme, literal string, line int) *Token {
	t := new(Token)
	t.Type = tokenType
	t.Lexeme = lexeme
	t.Literal = literal
	t.Line = line
	return t
}

// String returns the string representation of a Token
func (t *Token) String() string {
	return fmt.Sprintf("%d  %s  %v", t.Type, t.Lexeme, t.Literal)
}

// String return the string representation of a TokenType
func (t TokenType) String() string {

	switch t {
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"
	case Comma:
		return ","
	case Dot:
		return "."
	case Minus:
		return "-"
	case Plus:
		return "+"
	case Semicolon:
		return ";"
	case Slash:
		return "/"
	case Star:
		return "*"
	case Bang:
		return "!"
	case BangEqual:
		return "!="
	case Equal:
		return "="
	case EqualEqual:
		return "=="
	case Greater:
		return ">"
	case GreaterEqual:
		return ">="
	case Less:
		return "<"
	case LessEqual:
		return "<="
	case Identifier:
		return "identifier"
	case String:
		return "string"
	case Number:
		return "number"
	case And:
		return "and"
	case Class:
		return "class"
	case Else:
		return "else"
	case False:
		return "false"
	case Fun:
		return "fun"
	case For:
		return "for"
	case If:
		return "if"
	case Nil:
		return "nil"
	case Or:
		return "or"
	case Print:
		return "print"
	case Return:
		return "return"
	case Super:
		return "super"
	case This:
		return "this"
	case True:
		return "true"
	case Var:
		return "var"
	case While:
		return "while"
	case End:
		return "end"
	default:
		return "N/A"
	}
}
