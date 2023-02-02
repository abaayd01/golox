package main

import (
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

type astStructDefinition struct {
	Name string
}

type astTypeDefinition struct {
	astStructDefinition
	FieldList string
}

func defineAst(outputDir string, baseName string, types []string) error {
	outputFile := fmt.Sprintf("%s/ast.go", outputDir)
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
	}

	_, err = f.WriteString("package main\n\n")
	if err != nil {
		return err
	}

	t, err := template.New("astStruct").Parse("type {{ .Name }} struct{}\n")
	if err != nil {
		return err
	}

	err = t.Execute(f, astStructDefinition{Name: baseName})

	for _, s := range types {
		className := strings.Trim(strings.Split(s, ":")[0], " ")
		fieldNames := strings.ReplaceAll(strings.Trim(strings.Split(s, ":")[1], " "), ",", "\n")
		err = defineType(f, className, fieldNames)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	// formatting the file
	unfmtted, err := os.ReadFile(outputFile)
	if err != nil {
		return err
	}

	fmtted, err := format.Source(unfmtted)

	// overwriting the file
	f, err = os.Create(outputFile)
	if err != nil {
		return err
	}
	_, err = f.Write(fmtted)

	return err
}

func defineType(f io.Writer, className string, fieldList string) error {
	t, err := template.New("astStruct").Parse("type {{ .Name }} struct {\n{{ .FieldList }}\n}\n")
	if err != nil {
		return err
	}

	return t.Execute(f,
		astTypeDefinition{
			astStructDefinition: astStructDefinition{Name: className},
			FieldList:           fieldList,
		},
	)
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
