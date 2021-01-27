# Lox Grammar

The Lox grammar is defined below:

```
expression = literal | unary | binary | grouping ;
literal = NUMBER | STRING | BOOLEAN | NIL;
grouping = "(" expression ")" ;
unary = ( "-" | "!" ) expression ;
binary = expression operator expression ;
operator = "=" | "!=" | "<" | "<=" | ">" | ">="
    | "+" | "-" | "*" | "/" ;
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
