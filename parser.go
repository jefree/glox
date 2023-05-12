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

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() Expr {
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		expr = BinaryExpr{expr, p.previous(), p.comparison()}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		expr = BinaryExpr{expr, p.previous(), p.term()}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		expr = BinaryExpr{expr, p.previous(), p.factor()}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(MINUS, PLUS) {
		expr = BinaryExpr{expr, p.previous(), p.unary()}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		return UnaryExpr{p.previous(), p.unary()}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return LiteralExpr{false}
	}

	if p.match(TRUE) {
		return LiteralExpr{true}
	}

	if p.match(NIL) {
		return LiteralExpr{nil}
	}

	if p.match(NUMBER, STRING) {
		return LiteralExpr{p.previous().Literal}
	}

	if p.match(RIGHT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")

		return GroupingExpr{expr}
	}

	return LiteralExpr{p.peek().Literal}
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

func (p *Parser) consume(kind TokenKind, message string) Token {
	if p.check(kind) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) error(token Token, message string) ParseError {
	error(token, message)
	return ParseError{}
}
