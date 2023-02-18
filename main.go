package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

func (l *Lox) error(token Token, message string) {
	if token.tokenType == EOF {
		l.report(token.line, " at end", message)
	} else {
		l.report(token.line, fmt.Sprintf(" at '%s'", token.lexeme), message)
	}
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

	parser := Parser{
		Lox:     l,
		Tokens:  tokens,
		current: 0,
	}
	expr, err := parser.Parse()
	if err != nil {
		return err
	}

	printer := AstPrinter{}
	fmt.Println(printer.Print(expr))
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
