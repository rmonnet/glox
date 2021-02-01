# Lox Grammar

The Lox grammar is defined below:

```BNF
program = declaration* EOF ;
declaration = varDeclStmt | statement ;
varDeclStmt = "var" IDENTIFIER ( "=" expression )? ";" ;
statement = printStmt | exprStmt | block ;
exprStmt = expression ";" ;
printStmt = "print" expression ";" ;
block = "{" declaration* "}" ;
expression = assignment ;
assignment = IDENTIFIER "=" assignment | equality ;
equality = comparison ( ("!=" | "==" ) comparison )* ;
comparison = term ( (">" | ">=" | "<" | "<=" ) term )* ;
term = factor ( ( "-" | "+" ) factor )* ;
factor = unary ( ( "/" | "*" ) unary )* ;
unary = ( "!" | "-" ) unary
    | primary ;
primary = NUMBER | STRING | BOOLEAN | NIL
    | "(" expression ")"  | IDENTIFIER ;

NUMBER = [0-9]+ ( "." [0-9]+ )?
STRING = "\"" ( . )* "\""
BOOLEAN = "true" | "false"
NIL = "nil"
IDENTIFIER = ( [a-z] [A-Z] "_" ) ( [a-z] [A-Z] [0-9] "_" )*
```

Precedence rules (lowest to highest):

| Name       | Operator  | Associate |
| ---------- | --------- | --------- |
| Assignment | =         | right     |
| Equality   | == !=     | left      |
| Comparison | > >= < <= | left      |
| Term       | - +       | left      |
| Factor     | / *       | left      |
| Unary      | ! -       | right     |
