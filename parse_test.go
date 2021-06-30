package gogo

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		program string
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{name: "variable assignment", args: args{`foo := "bar"`}, want: []Token{
			{Identifier, "foo"},
			{Type: Whitespace},
			{Type: Assignment},
			{Type: Whitespace},
			{Type: OpenQuote},
			{Literal, "bar"},
			{Type: CloseQuote},
		}},
		{"variable reassignment", args{`foo := "bar"; foo = "bux"`}, []Token{
			{Identifier, "foo"},
			{Type: Whitespace},
			{Type: Assignment},
			{Type: Whitespace},
			{Type: OpenQuote},
			{Literal, "bar"},
			{Type: CloseQuote},
			{Type: Semicolon},
			{Type: Whitespace},
			{Identifier, "foo"},
			{Type: Whitespace},
			{Type: Reassignment},
			{Type: Whitespace},
			{Type: OpenQuote},
			{Literal, "bux"},
			{Type: CloseQuote},
		}},
		{"comments", args{"# hello!\nabc := 123"}, []Token{
			{Type: Newline},
			{Identifier, "abc"},
			{Type: Whitespace},
			{Type: Assignment},
			{Type: Whitespace},
			{Literal, "123"},
		}},
		{"function", args{"input()"}, []Token{
			{Identifier, "input"},
			{Type: OpenParen},
			{Type: CloseParen},
		}},
		{"function with string literal", args{`print("Hello, World!")`}, []Token{
			{Identifier, "print"},
			{Type: OpenParen},
			{Type: OpenQuote},
			{Literal, "Hello, World!"},
			{Type: CloseQuote},
			{Type: CloseParen},
		}},
		{"boolean literal", args{"x := true"}, []Token{
			{Identifier, "x"},
			{Type: Whitespace},
			{Type: Assignment},
			{Type: Whitespace},
			{Literal, "true"},
		}},
		{"shell command", args{"!ls -al"}, []Token{{ShellCmd, "ls -al"}}},
		{"shell command + boolean literal", args{"!ls -al\nx := true"}, []Token{
			{ShellCmd, "ls -al"},
			{Type: Newline},
			{Identifier, "x"},
			{Type: Whitespace},
			{Type: Assignment},
			{Type: Whitespace},
			{Literal, "true"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.program); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
