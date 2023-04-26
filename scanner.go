package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Scanner struct {
	Source  string
	Current int
	Start   int
	Line    int
	Tokens  []Token
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.Start = s.Current // set the start position for next token
		s.scanToken()
	}

	return s.Tokens
}

func (s *Scanner) scanToken() {
	ch := s.advance()

	switch ch {
	// single char tokens
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '/':
		s.addToken(SLASH, nil)
	case '*':
		s.addToken(STAR, nil)

	// one or two charecters tokens
	case '!':
		s.addMatchToken('=', BANG_EQUAL, BANG)
	case '=':
		s.addMatchToken('=', EQUAL_EQUAL, EQUAL)
	case '>':
		s.addMatchToken('=', GREATER_EQUAL, GREATER)
	case '<':
		s.addMatchToken('=', LESS_EQUAL, LESS)

	// ignore blank spaces
	case ' ':
	case '\r':
	case '\t':

	// new line
	case '\n':
		s.Line++

	// literal string
	case '"':
		s.string()

	// literal number, keywords and unknown tokens
	default:
		if isDigit(ch) {
			s.number()
		} else if isAlpha(ch) {
			s.identifier()

		} else {
			report(s.Line, "", fmt.Sprintf("Unexpected character '%c'.", ch))
		}
	}
}

func (s *Scanner) addToken(kind TokenKind, literal interface{}) {
	token := Token{
		kind,
		s.Source[s.Start:s.Current],
		literal,
		s.Line,
	}

	s.Tokens = append(s.Tokens, token)
}

func (s *Scanner) addMatchToken(ch byte, matchKind TokenKind, defaultKind TokenKind) {
	if s.match(ch) {
		s.addToken(matchKind, nil)
	} else {
		s.addToken(defaultKind, nil)
	}
}

func (s *Scanner) advance() byte {
	ch := s.Source[s.Current]
	s.Current++
	return ch
}

func (s *Scanner) match(ch byte) bool {
	if s.isAtEnd() {
		return false
	}

	next := s.Source[s.Current+1]

	if next != ch {
		return false
	}

	s.Current++
	return true
}

func (s *Scanner) peek() byte {
	return s.Source[s.Current]
}

func (s *Scanner) peekNext() byte {
	if s.Current+1 >= len(s.Source) {
		return '\x00'
	}

	return s.Source[s.Current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		report(s.Line, "", "Unterminated string.")
		return
	}

	s.advance()

	value := s.Source[s.Start+1 : s.Current-1]
	s.addToken(STRING, value)

}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	number, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)

	if err != nil {
		fmt.Println(err)
		os.Exit(SoftwareExitCode)
	}

	s.addToken(NUMBER, number)
}

func (s *Scanner) identifier() {
	for isAlphaDigit(s.peek()) {
		s.advance()
	}

	if tokenType, ok := keywords[s.Source[s.Start:s.Current]]; ok {
		s.addToken(tokenType, nil)
	} else {
		s.addToken(IDENTIFIER, nil)
	}
}

func isAlpha(ch byte) bool {
	charRune := rune(ch)
	return unicode.IsLetter(charRune) || charRune == '_'
}

func isDigit(ch byte) bool {
	return unicode.IsNumber(rune(ch))
}

func isAlphaDigit(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Source: source,
	}
}
