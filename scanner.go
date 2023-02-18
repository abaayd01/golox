package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// Object - see how this needs to be used later
type Object interface{}

type TokenType int

const (
	// single character tokens

	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	EOF
)

func (t TokenType) String() string {
	switch t {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FUN:
		return "FUN"
	case FOR:
		return "FOR"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case EOF:
		return "EOF"
	default:
		return ""
	}
}

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   Object
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", t.tokenType, t.lexeme, t.literal)
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	lox     *Lox
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *Scanner) scanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens)
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	r := s.advance()
	switch r {
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		s.line++
		break
	case '(':
		s.addToken(LEFT_PAREN)
		break
	case ')':
		s.addToken(RIGHT_PAREN)
		break
	case '{':
		s.addToken(LEFT_BRACE)
		break
	case '}':
		s.addToken(RIGHT_BRACE)
		break
	case ',':
		s.addToken(COMMA)
		break
	case '.':
		s.addToken(DOT)
		break
	case '-':
		s.addToken(MINUS)
		break
	case '+':
		s.addToken(PLUS)
		break
	case ';':
		s.addToken(SEMICOLON)
		break
	case '*':
		s.addToken(STAR)
		break
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
		break
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
		break
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
		break
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
		break
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.scanMultiLineComment()
		} else {
			s.addToken(SLASH)
		}
	case '"':
		s.scanStringLiteral()
		break
	default:
		if isDigit(r) {
			s.scanNumberLiteral()
		} else if isAlpha(r) {
			s.scanIdentifierOrKeyword()
		} else {
			s.lox.reportError(s.line, fmt.Sprintf("Unexpected character: %s.", string(r)))
		}
		break
	}
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isAlpha(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func (s *Scanner) scanStringLiteral() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.lox.reportError(s.line, "Unterminated string.")
		return
	}

	s.advance()

	strVal := s.source[s.start+1 : s.current-1]
	s.addLiteralToken(STRING, strVal)
}

func (s *Scanner) scanNumberLiteral() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	val, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.lox.reportError(s.line, fmt.Sprintf("Could not parse as float: %s.", s.source[s.start:s.current]))
	}
	s.addLiteralToken(NUMBER, val)
}

func (s *Scanner) scanIdentifierOrKeyword() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	val := s.source[s.start:s.current]
	kw, ok := keywords[val]

	if !ok {
		s.addToken(IDENTIFIER)
		return
	}

	s.addToken(kw)
}

func (s *Scanner) scanMultiLineComment() {
	for s.peek() != '*' && s.peekNext() != '/' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.lox.reportError(s.line, "Unterminated multi-line comment.")
		return
	}

	// need to advance twice at end to consume */ closing tag
	s.advance()
	s.advance()
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return rune(r)
}

func (s *Scanner) addToken(t TokenType) {
	txt := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		tokenType: t,
		lexeme:    txt,
		literal:   nil,
		line:      s.line,
	})
}

func (s *Scanner) addLiteralToken(t TokenType, val interface{}) {
	txt := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		tokenType: t,
		lexeme:    txt,
		literal:   val,
		line:      s.line,
	})
}

func (s *Scanner) match(r rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != r {
		return false
	}

	s.advance()
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}

	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return rune(s.source[s.current+1])
}
