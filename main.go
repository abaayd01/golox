package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Lox struct {
	hadError        bool
	hadRuntimeError bool
}

func (l *Lox) reportError(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	_, _ = fmt.Fprintf(os.Stderr, "[line: %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}

func (l *Lox) error(token Token, message string) {
	if token.tokenType == EOF {
		l.report(token.line, " at end", message)
	} else {
		l.report(token.line, fmt.Sprintf(" at '%s'", token.lexeme), message)
	}
}

func (l *Lox) runtimeError(err RuntimeError) {
	_, _ = fmt.Fprintf(os.Stderr, "[line %d] %s\n", err.Token.line, err.Error())
	l.hadRuntimeError = true
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
	statements, err := parser.Parse()
	if err != nil { // stop if there was a parsing error
		return err
	}

	e := &Environment{Values: map[string]any{}}
	interpreter := Interpreter{
		Lox:         l,
		Environment: e,
	}

	_ = interpreter.InterpretStatements(statements)
	//_, _ = interpreter.InterpretExpression(expr) // don't blow up if there's runtime errors?

	// temporary AstPrinter code
	//printer := AstPrinter{}
	//fmt.Println(printer.Print(expr))
	return nil
}

func (l *Lox) runFile(path string) error {
	log.Printf("runFile, path: %s", path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("runFile error, os.ReadFile: %w", err)
	}

	// TODO: re-think how errors are propagated up from Parser + Interpreter into Lox entry point.
	// Need to figure out how to best deal with runtime errors vs. parsing errors.
	// At the moment run doesn't return an error for runtime errors, the Interpreter will just set the hadRuntimeError flag.
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
			log.Fatalln(err)
		}

		return
	}

	_ = l.runPrompt()
}
