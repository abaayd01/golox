package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

type Lox struct {
	hadError bool
}

func (l *Lox) reportError(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	_, _ = fmt.Fprintf(os.Stderr, "[line: %d] Error%s: %s]\n", line, where, message)
	l.hadError = true
}

func (l *Lox) run(source string) error {
	scanner := Scanner{
		lox:    l,
		source: source,
	}
	tokens, err := scanner.scanTokens()
	if err != nil {
		return err
	}

	// just print out the tokens for now
	for _, token := range tokens {
		fmt.Println(token, token.line)
	}

	return nil
}

func (l *Lox) runFile(path string) error {
	log.Printf("runFile, path: %s", path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("runFile error, os.ReadFile: %w", err)
	}

	return l.run(string(bytes))
}

func (l *Lox) printPrompt() {
	fmt.Printf("> ")
}

func (l *Lox) runPrompt() error {
	reader := bufio.NewScanner(os.Stdin)
	l.printPrompt()
	for reader.Scan() {
		// just echoing out the input for now
		cmd := reader.Text()
		fmt.Println(cmd)
		err := l.run(cmd)
		if err != nil {
			return err
		}
		l.printPrompt()
	}
	// Print an additional line if we encountered an EOF character
	fmt.Println()
	return nil
}

func main() {
	args := os.Args

	if len(args) > 2 {
		log.Fatal("Invalid usage")
	}

	l := Lox{}

	if len(args) == 2 {
		err := l.runFile(args[1])
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	_ = l.runPrompt()
}
