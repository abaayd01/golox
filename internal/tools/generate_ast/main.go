package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

func defineAst(outputDir string, baseName string, types []string) error {
	outputFile := fmt.Sprintf("%s/%s.go", outputDir, strings.ToLower(baseName))
	buf := bytes.Buffer{}

	_, err := buf.WriteString(`
	// this file is auto-generated with bin/generate_ast
	// DO NOT EDIT
	package main

	`)
	if err != nil {
		return fmt.Errorf("error writing string to buf: %w", err)
	}

	tmpl, err := template.New("exprStruct").Parse(
		`type {{.Name}} interface {
			Accept(visitor {{.Name}}Visitor) (any, error)
		}
		`,
	)
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, struct{ Name string }{Name: baseName})
	if err != nil {
		return err
	}

	for _, s := range types {
		className := strings.Trim(strings.Split(s, ":")[0], " ")
		fieldNames := strings.Split(strings.Trim(strings.Split(s, ":")[1], " "), ",")
		err = defineType(&buf, className, fieldNames, baseName)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	err = defineVisitor(&buf, types, baseName)
	if err != nil {
		return err
	}

	// formatting the temporary buffer
	fmtted, err := format.Source(buf.Bytes())

	fmt.Printf("buf: %s", string(buf.Bytes()))
	fmt.Printf("fmtted: %s", fmtted)

	// writing it out to file
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error writing buf to outputFile: %w", err)
	}
	_, err = f.Write(fmtted)

	return err
}

func defineType(w io.Writer, typeName string, fieldList []string, baseName string) error {
	t, err := template.New("astStruct").Parse(
		`type {{ .Name }} struct {
			{{range .FieldList}}
			{{.}}{{end}}
		}
		func (t {{.Name}}) Accept(visitor {{.BaseName}}Visitor) (any, error) {
			return visitor.Visit{{.Name}}(t)
		}
		`,
	)
	if err != nil {
		return err
	}

	return t.Execute(w,
		struct {
			Name      string
			FieldList []string
			BaseName  string
		}{
			Name:      typeName,
			FieldList: fieldList,
			BaseName:  baseName,
		},
	)
}

func defineVisitor(w io.Writer, types []string, baseName string) error {
	var typeNames []string
	for _, s := range types {
		typeName := strings.Trim(strings.Split(s, ":")[0], " ")
		typeNames = append(typeNames, typeName)
	}

	tmpl, err := template.New("visitorInterface").Parse(
		`type {{.BaseName}}Visitor interface {
				{{range .TypeNames}}
				Visit{{.}}(expr {{.}}) (any, error){{end}}
			}
		`,
	)
	if err != nil {
		return fmt.Errorf("could not parse visitorInterface template: %w", err)
	}

	err = tmpl.Execute(w, struct {
		BaseName  string
		TypeNames []string
	}{
		BaseName:  baseName,
		TypeNames: typeNames,
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	args := os.Args

	if len(args) != 2 {
		log.Fatal("Invalid usage...")
	}

	outputDir := args[1]

	err := defineAst(outputDir, "Expr", []string{
		"Unary: operator Token, right Expr",
		"Binary: left Expr, operator Token, right Expr",
		"Grouping: expression Expr",
		"Literal: value Object",
		"Var: name Token",
		"Assign: name Token, value Expr",
		"Logical: left Expr, operator Token, right Expr",
	})
	if err != nil {
		log.Fatal(err)
	}

	err = defineAst(outputDir, "Stmt", []string{
		"StmtExpression: expression Expr",
		"StmtPrint: expression Expr",
		"StmtVar: name Token, initializer Expr",
		"StmtBlock: statements []Stmt",
		"StmtIf: condition Expr, thenBranch Stmt, elseBranch Stmt",
		"StmtWhile: condition Expr, body Stmt",
	})

	if err != nil {
		log.Fatal(err)
	}
}
