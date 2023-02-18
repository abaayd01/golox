package main

import (
	"fmt"
	"strconv"
	"unicode"
)

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
