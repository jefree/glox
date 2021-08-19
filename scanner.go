package main

import "fmt"

type Scanner struct {
	source  string
	current int
	start   int
	line    int
	tokens  []Token
}

func (s *Scanner) ScanTokens() []Token {
	for s.current < len(s.source) {
		s.start = s.current // set the start position for next token
		s.scanToken()
	}

	return s.tokens
}

func (s *Scanner) scanToken() {
	ch := s.source[s.current]

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
		s.addMatchToken('=', EQUAL, EQUAL_EQUAL)
	case '>':
		s.addMatchToken('=', GREATER, GREATER_EQUAL)
	case '<':
		s.addMatchToken('=', LESS, LESS_EQUAL)
		// literal, keywords and unknown tokens
	case '\n':
		s.line++
	default:
		fmt.Println("unknown char: ", string(ch))
	}

	s.current++
}

func (s *Scanner) addToken(kind TokenKind, literal interface{}) {
	token := Token{
		kind,
		s.source[s.start : s.current+1],
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

func (s *Scanner) match(ch byte) bool {
	next := s.source[s.current+1]

	if next != ch {
		return false
	}

	s.current++
	return true
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
	}
}
