# Lox Grammar

The Lox grammar is defined below:

```BNF
expression = equality ;
equality = comparison ( ("!=" | "==" ) comparison )* ;
comparison = term ( (">" | ">=" | "<" | "<=" ) term )* ;
term = factor ( ( "-" | "+" ) factor )* ;
factor = unary ( ( "/" | "*" ) unary )* ;
unary = ( "!" | "-" ) unary
    | primary ;
primary = NUMBER | STRING | BOOLEAN | NIL
    | "(" expression ")" ;

NUMBER = [0-9]+ ( "." [0-9]+ )?
STRING = "\"" ( . )* "\""
BOOLEAN = "true" | "false"
NIL = "nil"
```

Precedence rules (lowest to highest):

| Name       | Operator  | Associate |
| ---------- | --------- | --------- |
| Equality   | == !=     | left      |
| Comparison | > >= < <= | left      |
| Term       | - +       | left      |
| Factor     | / *       | left      |
| Unary      | ! -       | right     |
