# Lox Grammar

The Lox grammar is defined below:

```BNF
program = 
    declaration* EOF ;

declaration = 
    classDeclStmt |funDeclStmt | varDeclStmt | statement ;

classDeclStmt =
    "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}" ;

funDeclStmt =
    "fun" function;

function =
    IDENTIFIER "(" parameters? ")" block ;

parameters =
    IDENTIFIER ( "," IDENTIFIER )* ;

varDeclStmt =
    "var" IDENTIFIER ( "=" expression )? ";" ;

statement =
    exprStmt | forStmt | ifStmt | printStmt | returnStmt 
    | whileStmt | block ;

exprStmt =
    expression ";" ;

forStmt =
    "for" "(" ( varDecl | exprStmt | ";" )
    expression? ";" expression? ")" statement ;

ifStmt =
    "if" "(" expression ")" statement ( "else" statement )? ;

printStmt =
    "print" expression ";" ;

returnStmt =
    "return" expression? ";" ;

whileStmt =
    "while" "(" expression ")" statement ;

block =
    "{" declaration* "}" ;

expression =
    assignment ;

assignment =
    ( call "." )? IDENTIFIER "=" assignment | logic_or ;

logic_or =
    logic_and ( "or" logic_and )* ;

logic_and =
    equality ( "and" equality )* ;

equality =
    comparison ( ("!=" | "==" ) comparison )* ;

comparison =
    term ( (">" | ">=" | "<" | "<=" ) term )* ;

term =
    factor ( ( "-" | "+" ) factor )* ;

factor =
    unary ( ( "/" | "*" ) unary )* ;

unary =
    ( "!" | "-" ) unary | call ;

call =
    primary ( "(" arguments? ")" | "." IDENTIFIER )* ;

arguments =
    expression ( "," expression )* ;

primary =
    NUMBER | STRING | BOOLEAN | NIL | "(" expression ")"
    | "this" | "super" | IDENTIFIER ;

NUMBER =
    [0-9]+ ( "." [0-9]+ )?

STRING =
    "\"" ( . )* "\""

BOOLEAN =
    "true" | "false"

NIL =
    "nil"

IDENTIFIER =
    ( [a-z] [A-Z] "_" ) ( [a-z] [A-Z] [0-9] "_" )*
```

Precedence rules (lowest to highest):

| Name       | Operator  | Associate |
| ---------- | --------- | --------- |
| Assignment | =         | right     |
| Or         | or        | left      |
| And        | and       | left      |
| Equality   | == !=     | left      |
| Comparison | > >= < <= | left      |
| Term       | - +       | left      |
| Factor     | / *       | left      |
| Unary      | ! -       | right     |
