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
	outputFile := fmt.Sprintf("%s/ast.go", outputDir)
	buf := bytes.Buffer{}

	_, err := buf.WriteString("package main\n\n")
	if err != nil {
		return err
	}

	tmpl, err := template.New("exprStruct").Parse(
		`type {{.Name}} interface {
			Accept(visitor Visitor) any
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
		err = defineType(&buf, className, fieldNames)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	err = defineVisitor(&buf, types)
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
		return err
	}
	_, err = f.Write(fmtted)

	return err
}

func defineType(w io.Writer, typeName string, fieldList []string) error {
	t, err := template.New("astStruct").Parse(
		`type {{ .Name }} struct {
			{{range .FieldList}}
			{{.}}{{end}}
		}
		func (t {{.Name}}) Accept(visitor Visitor) any {
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
		}{
			Name:      typeName,
			FieldList: fieldList,
		},
	)
}

func defineVisitor(w io.Writer, types []string) error {
	var typeNames []string
	for _, s := range types {
		typeName := strings.Trim(strings.Split(s, ":")[0], " ")
		typeNames = append(typeNames, typeName)
	}

	tmpl, err := template.New("visitorInterface").Parse(
		`type Visitor interface {
				{{range .TypeNames}}
				Visit{{.}}(expr {{.}}) any{{end}}
			}
		`,
	)
	if err != nil {
		return fmt.Errorf("could not parse visitorInterface template: %w", err)
	}

	err = tmpl.Execute(w, struct {
		TypeNames []string
	}{
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
	})
	if err != nil {
		log.Fatal(err)
	}
}
