package main

// expression → equality ;
// equality   → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term       → factor ( ( "-" | "+" ) factor )* ;
// factor     → unary ( ( "/" | "*" ) unary )* ;
// unary      → ( "!" | "-" ) unary
//            | primary ;
// primary    → NUMBER | STRING | "true" | "false" | "nil"
//            | "(" expression ")" ;

type ParseError struct{}

func (err ParseError) Error() string {
	return "ParseError"
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{expr, p.previous(), right}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{expr, p.previous(), right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{expr, p.previous(), right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{expr, p.previous(), right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		expr, err := p.unary()
		if err != nil {
			return nil, err
		}

		return UnaryExpr{p.previous(), expr}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return LiteralExpr{false}, nil
	}

	if p.match(TRUE) {
		return LiteralExpr{true}, nil
	}

	if p.match(NIL) {
		return LiteralExpr{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return LiteralExpr{p.previous().Literal}, nil
	}

	if p.match(RIGHT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return GroupingExpr{expr}, nil
	}

	return nil, ParseError{}
}

func (p *Parser) match(kinds ...TokenKind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(kind TokenKind, message string) error {
	if p.check(kind) {
		p.advance()
		return nil
	}

	return p.fail(p.peek(), message)
}

func (p *Parser) fail(token Token, message string) ParseError {
	fail(token, message)

	return ParseError{}
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() Token {
	ch := p.tokens[p.current]
	p.current++
	return ch
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) check(kind TokenKind) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Kind == kind
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == EOF
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind == SEMICOLON {
			return
		}

		switch p.peek().Kind {
		case CLASS, FOR, FUN, IF, PRINT, RETURN, VAR, WHILE:
			return
		default:
			p.advance()
		}
	}
}
