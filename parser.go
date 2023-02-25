package main

import (
	"errors"
)

type Parser struct {
	Lox     *Lox
	Tokens  []Token
	current int
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

// Grammar Production:
// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

// Grammar Production:
// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

// Grammar Production:
// term → factor ( ( "-" | "+" ) factor )*;
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

// Grammar Production:
// factor → unary ( ( "/" | "*" ) unary )*;
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

// Grammar Production:
// unary → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{
			operator: op,
			right:    right,
		}, nil
	}

	return p.primary()
}

// Grammar Production:
// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() (Expr, error) {
	if p.match(NUMBER, STRING) {
		return Literal{value: p.previous().literal}, nil
	}

	if p.match(TRUE) {
		return Literal{value: true}, nil
	}

	if p.match(FALSE) {
		return Literal{value: false}, nil
	}

	if p.match(NIL) {
		return Literal{value: nil}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return Grouping{expression: expr}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) peek() Token {
	return p.Tokens[p.current]
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) previous() Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) consume(tokenType TokenType, message string) (*Token, error) {
	if p.check(tokenType) {
		t := p.advance()
		return &t, nil
	}

	return nil, p.error(p.peek(), message)
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

// error handling code
var ParseError = errors.New("parse error")

func (p *Parser) error(token Token, message string) error {
	p.Lox.error(token, message)
	return ParseError
}

func (p *Parser) isAtEndOfPreviousStatement() bool {
	return p.previous().tokenType == SEMICOLON
}

// map of keywords which start a statement
var statementStarterKeywords = map[TokenType]bool{
	CLASS:  true,
	FUN:    true,
	VAR:    true,
	FOR:    true,
	IF:     true,
	WHILE:  true,
	PRINT:  true,
	RETURN: true,
}

func (p *Parser) isAtStartOfNewStatement() bool {
	_, ok := statementStarterKeywords[p.peek().tokenType]
	return ok
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.isAtEndOfPreviousStatement() || p.isAtStartOfNewStatement() {
			return
		}
		p.advance()
	}
}
