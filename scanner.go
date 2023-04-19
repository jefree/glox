package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Scanner struct {
	source  string
	current int
	start   int
	line    int
	tokens  []Token
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current // set the start position for next token
		s.scanToken()
	}

	return s.tokens
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
		s.line++

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
			report(s.line, "", fmt.Sprintf("Unexpected character '%c'.", ch))
		}
	}
}

func (s *Scanner) addToken(kind TokenKind, literal interface{}) {
	token := Token{
		kind,
		s.source[s.start:s.current],
		literal,
		s.line,
	}

	s.tokens = append(s.tokens, token)
}

func (s *Scanner) addMatchToken(ch byte, matchKind TokenKind, defaultKind TokenKind) {
	if s.match(ch) {
		s.addToken(matchKind, nil)
	} else {
		s.addToken(defaultKind, nil)
	}
}

func (s *Scanner) advance() byte {
	ch := s.source[s.current]
	s.current++
	return ch
}

func (s *Scanner) match(ch byte) bool {
	if s.isAtEnd() {
		return false
	}

	next := s.source[s.current+1]

	if next != ch {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() byte {
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}

	return s.source[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		report(s.line, "", "Unterminated string.")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
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

	number, err := strconv.ParseFloat(s.source[s.start:s.current], 64)

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

	if tokenType, ok := keywords[s.source[s.start:s.current]]; ok {
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
		source: source,
	}
}
