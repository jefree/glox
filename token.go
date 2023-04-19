package main

import "fmt"

type Token struct {
	kind    TokenKind
	lexeme  string
	literal interface{}
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s: '%s'\tline: %d", t.kind, t.lexeme, t.line)
}
