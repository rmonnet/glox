
# What is the glox project?

This project implements the code from the book [crafting interpreters](https://craftinginterpreters.com) in go.

The first version literally follows the book and doesn't attempt to use go specific features.

The lox grammar is defined [here](grammar.md).

The `lox/lang` package includes the AST (ast.go),
the tokens (tokens.go), scanner (scanner.go)
and parser (parser.go). The scanner takes a script
as `string` and returns a slice of tokens (`[]*Tokens`).
The parser takes a slice of tokens and returns a slice of AST
nodes (`[]Stmt`).

The `lox/interp` package includes the interpreter itself (interp.go) and the resolver (resolver.go). The resolver
performs static analysis. It could be seen as being part
of the `lang` package since it checks for compile errors 
but it is bundled with the interpreter in the original text.
It also has one direct call to the interpreter (`Interp.Resolve()`) which makes it dependent on the `interp` package.

There are unit tests for the low level `lang` package
and the interpreter itself. The interpreter tests are
written as go testable example since it makes it very 
readable. Each test consists of a lox script and the 
expected results as `// Ouput:` comments.

The main deviations from the original java implementation
described in the [crafting interpreters](https://craftinginterpreters.com) are:

- The AST statements (`lang.Stmt`) and expressions (`lang.Expr`) do not use the visitor pattern as in the java code.
- The java code uses `Object` for the dynamic values in expressions. The go code uses `interface{}`
- The AST Nodes implement `PrettyPrint()` which allows to pretty print any AST tree without special package. 

# Installation

You can install the lox interpreter using `go get gitlab.com/rcmonnet/glox`.

# Example

See the examples in the `examples` directory. They are
mostly taken from the original text.

# Build the project

You can also install the source code on your machine by typing:

```
git clone gitlab.com/rcmonnet/glox`
cd glox
go install gitlab.com/rcmonnet/glox
```

# FAQ

# Changes

# License

This is an implementation of the lox interpreter presented
by Bob Nystrom [here](https://craftinginterpreters.com/).

The go code in this project is licensed under the [MIT
license](https://mit-license.org/).
