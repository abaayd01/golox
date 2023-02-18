package main

import (
	"testing"
)

func TestAstPrinter_Print(t *testing.T) {
	type args struct {
		expr Expr
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				expr: Binary{
					left: Literal{value: 1},
					operator: Token{
						tokenType: PLUS,
						lexeme:    "+",
						literal:   nil,
						line:      0,
					},
					right: Literal{value: 2},
				},
			},
			want: "(+ 1 2)",
		},
		{
			name: "complex",
			args: args{
				expr: Binary{
					left: Unary{
						operator: Token{
							tokenType: MINUS,
							lexeme:    "-",
							literal:   nil,
							line:      0,
						},
						right: Literal{value: 123},
					},
					operator: Token{
						tokenType: STAR,
						lexeme:    "*",
						literal:   nil,
						line:      0,
					},
					right: Grouping{expression: Literal{value: 45.67}},
				},
			},
			want: "(* (- 123) (group 45.67))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AstPrinter{}
			if got := a.Print(tt.args.expr); got != tt.want {
				t.Errorf("Print() = %v, want %v", got, tt.want)
			}
		})
	}
}
