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
		{"variable assignment", args{`foo := "bar"`}, []Token{
			{Identifier, "foo"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
			{OpenQuote, ""},
			{Literal, "bar"},
			{CloseQuote, ""},
		}},
		{"variable reassignment", args{`foo := "bar"; foo = "bux"`}, []Token{
			{Identifier, "foo"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
			{OpenQuote, ""},
			{Literal, "bar"},
			{CloseQuote, ""},
			{Semicolon, ""},
			{Whitespace, ""},
			{Identifier, "foo"},
			{Whitespace, ""},
			{Reassignment, ""},
			{Whitespace, ""},
			{OpenQuote, ""},
			{Literal, "bux"},
			{CloseQuote, ""},
		}},
		{"comments", args{"# hello!\nabc := 123"}, []Token{
			{Newline, ""},
			{Identifier, "abc"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
			{Literal, "123"},
		}},
		{"function", args{"input()"}, []Token{
			{Identifier, "input"},
			{OpenParen, ""},
			{CloseParen, ""},
		}},
		{"function with string literal", args{`print("Hello, World!")`}, []Token{
			{Identifier, "print"},
			{OpenParen, ""},
			{OpenQuote, ""},
			{Literal, "Hello, World!"},
			{CloseQuote, ""},
			{CloseParen, ""},
		}},
		{"boolean literal", args{"x := true"}, []Token{
			{Identifier, "x"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
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
