package main

import "fmt"

type Token struct {
	kind    TokenKind
	lexeme  string
	literal interface{}
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("char: '%s'\tline: %d", t.lexeme, t.line)
}
