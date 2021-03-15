// Package lang contains definitions for the lox
// language grammar.
// That includes tokens, AST Nodes, the scanner to generate
// tokens and the parser to generate an AST tree.
package lang

import (
	"fmt"
	"strings"
)

// TokenType represents the type of a lox token.
type TokenType int

const (
	// EndToken is a special token that represents the end of stream.
	EndToken TokenType = iota
	// AndToken represents an 'and' token.
	AndToken
	// BangToken represents a '!' token.
	BangToken
	// BangEqualToken represents a '!=' token.
	BangEqualToken
	// ClassToken represents a 'class' token.
	ClassToken
	// CommaToken represents a ',' token.
	CommaToken
	// DotToken represents a '.' token.
	DotToken
	// ElseToken represents an 'else' token.
	ElseToken
	// EqualToken represents an '=' token.
	EqualToken
	// EqualEqualToken represents an '==' token.
	EqualEqualToken
	// FalseToken represents a 'false' token.
	FalseToken
	// FunToken represents a 'fun' token.
	FunToken
	// ForToken represents a 'for' token.
	ForToken
	// GreaterToken represents a '>' token.
	GreaterToken
	// GreaterEqualToken represents a '>=' token.
	GreaterEqualToken
	// IdentifierToken represents any identifier token.
	IdentifierToken
	// IfToken represents an 'if' token.
	IfToken
	// LeftBraceToken represents a '{' token.
	LeftBraceToken
	// LeftParenToken represents a '(' token.
	LeftParenToken
	// LessToken represents a '<'' token.
	LessToken
	// LessEqualToken represents a '<=' token.
	LessEqualToken
	// MinusToken represents a '-' token.
	MinusToken
	// NilToken represents a 'nil' token.
	NilToken
	// NumberToken represents any number token.
	NumberToken
	// OrToken represents an 'or' token.
	OrToken
	// PlusToken represents a '+' token.
	PlusToken
	// PrintToken represents a 'print' token.
	PrintToken
	// ReturnToken represents a 'return' token.
	ReturnToken
	// RightBraceToken represents a '}' token.
	RightBraceToken
	// RightParenToken represents a ')' token.
	RightParenToken
	// SemicolonToken represents a ';' token.
	SemicolonToken
	// SlashToken represents a '/' token.
	SlashToken
	// StarToken represents a '*' token.
	StarToken
	// StringToken represents any string token.
	StringToken
	// SuperToken represents a 'super' token.
	SuperToken
	// ThisToken represents a 'this' token.
	ThisToken
	// TrueToken represents a 'true' token.
	TrueToken
	// VarToken represents a 'var' token.
	VarToken
	// WhileToken represents a 'while' token.
	WhileToken
)

// Token represents a lox token.
type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

// String returns the string representation of a Token.
func (t *Token) String() string {

	switch t.Type {
	case IdentifierToken:
		return fmt.Sprintf("Identifier(%s)", t.Lexeme)
	case NumberToken:
		return fmt.Sprintf("Number(%s)", t.Lexeme)
	case StringToken:
		value := strings.Trim(t.Lexeme, "\"")
		return fmt.Sprintf("String(%s)", value)
	default:
		return t.Type.String()
	}
}

// String return the string representation of a TokenType.
func (t TokenType) String() string {

	switch t {
	case EndToken:
		return "end-of-stream"
	case AndToken:
		return "and"
	case BangToken:
		return "!"
	case BangEqualToken:
		return "!="
	case ClassToken:
		return "class"
	case CommaToken:
		return ","
	case DotToken:
		return "."
	case ElseToken:
		return "else"
	case EqualToken:
		return "="
	case EqualEqualToken:
		return "=="
	case FalseToken:
		return "false"
	case FunToken:
		return "fun"
	case ForToken:
		return "for"
	case GreaterToken:
		return ">"
	case GreaterEqualToken:
		return ">="
	case IdentifierToken:
		return "identifier"
	case IfToken:
		return "if"
	case LeftBraceToken:
		return "{"
	case LeftParenToken:
		return "("
	case LessToken:
		return "<"
	case LessEqualToken:
		return "<="
	case MinusToken:
		return "-"
	case NilToken:
		return "nil"
	case NumberToken:
		return "number"
	case PlusToken:
		return "+"
	case RightParenToken:
		return ")"
	case RightBraceToken:
		return "}"
	case SemicolonToken:
		return ";"
	case SlashToken:
		return "/"
	case StarToken:
		return "*"
	case StringToken:
		return "string"
	case OrToken:
		return "or"
	case PrintToken:
		return "print"
	case ReturnToken:
		return "return"
	case SuperToken:
		return "super"
	case ThisToken:
		return "this"
	case TrueToken:
		return "true"
	case VarToken:
		return "var"
	case WhileToken:
		return "while"
	default:
		return "invalid-token"
	}
}
