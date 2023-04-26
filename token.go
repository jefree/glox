package main

import "fmt"

type Token struct {
	Kind    TokenKind
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s: '%s'\tline: %d", t.Kind, t.Lexeme, t.Line)
}
