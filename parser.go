package main

import (
	"errors"
	"fmt"
	"reflect"
)

type Parser struct {
	Lox     *Lox
	Tokens  []Token
	current int
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt

	for !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}

	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(FUN) {
		return p.function("function")
	}

	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) function(kind string) (Stmt, error) {
	name, err := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s %s", kind, name))
	if err != nil {
		return nil, err
	}

	var parameters []Token

	if !p.check(RIGHT_PAREN) {
		newParam, err := p.consume(IDENTIFIER, "Expect parameter name.")
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, *newParam)

		for p.match(COMMA) {
			if len(parameters) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 parameters.")
			}

			newParam, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, *newParam)
		}
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return StmtFunction{
		name:   *name,
		params: parameters,
		body:   body.(StmtBlock),
	}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return StmtVar{
		name:        *name,
		initializer: initializer,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(RETURN) {
		return p.returnStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(LEFT_BRACE) {
		return p.blockStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		v, err := p.varDeclaration()
		if err != nil {
			return nil, err
		}

		initializer = v
	} else {
		v, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}

		initializer = v
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		v, err := p.expression()
		if err != nil {
			return nil, err
		}

		condition = v
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after loop condition.")
	var increment Expr
	if !p.check(RIGHT_PAREN) {
		v, err := p.expression()
		if err != nil {
			return nil, err
		}

		increment = v
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after 'for' clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// transforming this into the AST for a while loop
	if increment != nil {
		body = StmtBlock{
			statements: []Stmt{
				body,
				StmtExpression{
					expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = Literal{
			value: true,
		}
	}

	body = StmtWhile{
		condition: condition,
		body:      body,
	}

	if initializer != nil {
		body = StmtBlock{
			statements: []Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after 'if' condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	if p.match(ELSE) {
		elseBranch, err := p.statement()
		if err != nil {
			return nil, err
		}

		return StmtIf{
			condition:  condition,
			thenBranch: thenBranch,
			elseBranch: elseBranch,
		}, nil
	}

	return StmtIf{
		condition:  condition,
		thenBranch: thenBranch,
	}, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return StmtPrint{expression: expr}, nil
}

func (p *Parser) returnStatement() (Stmt, error) {
	returnKeyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consume(SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return StmtReturn{
		returnKeyword: returnKeyword,
		value:         value,
	}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return StmtExpression{expression: expr}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after 'while' condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	fmt.Println(body)
	return StmtWhile{
		condition: condition,
		body:      body,
	}, nil
}

func (p *Parser) blockStatement() (Stmt, error) {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		s, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, s)
	}

	_, err := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return StmtBlock{
		statements: statements,
	}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.logicalOr() // lhs of eq
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		eq := p.previous()
		value, err := p.assignment() // evaluating the rhs
		if err != nil {
			return nil, err
		}

		isVar := reflect.TypeOf(expr).String() == "main.Var"

		if isVar {
			name := expr.(Var).name
			return Assign{
				name:  name,
				value: value,
			}, nil
		}

		return nil, p.error(eq, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) logicalOr() (Expr, error) {
	expr, err := p.logicalAnd()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		op := p.previous()
		right, err := p.logicalAnd()
		if err != nil {
			return nil, err
		}

		return Logical{
			left:     expr,
			right:    right,
			operator: op,
		}, nil

	}

	return expr, nil
}

func (p *Parser) logicalAnd() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		return Logical{
			left:     expr,
			right:    right,
			operator: op,
		}, nil

	}

	return expr, nil
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

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(expr Expr) (Expr, error) {
	var arguments []Expr

	if !p.check(RIGHT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, expr)

		for p.match(COMMA) {
			if len(arguments) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}

			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			arguments = append(arguments, expr)
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return Call{
		callee:    expr,
		arguments: arguments,
		paren:     *paren,
	}, nil
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

	if p.match(IDENTIFIER) {
		return Var{name: p.previous()}, nil
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
