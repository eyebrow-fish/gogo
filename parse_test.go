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
			{Literal, "foo"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
			{OpenQuote, ""},
			{Literal, "bar"},
			{CloseQuote, ""},
		}},
		{"variable reassignment", args{`foo := "bar"; foo = "bux"`}, []Token{
			{Literal, "foo"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
			{OpenQuote, ""},
			{Literal, "bar"},
			{CloseQuote, ""},
			{Semicolon, ""},
			{Whitespace, ""},
			{Literal, "foo"},
			{Whitespace, ""},
			{Reassignment, ""},
			{Whitespace, ""},
			{OpenQuote, ""},
			{Literal, "bux"},
			{CloseQuote, ""},
		}},
		{"comments", args{"# hello!\nabc := 123"}, []Token{
			{Newline, ""},
			{Literal, "abc"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
			{Literal, "123"},
		}},
		{"function", args{"input()"}, []Token{
			{Literal, "input"},
			{OpenParen, ""},
			{CloseParen, ""},
		}},
		{"function with string literal", args{`print("Hello, World!")`}, []Token{
			{Literal, "print"},
			{OpenParen, ""},
			{OpenQuote, ""},
			{Literal, "Hello, World!"},
			{CloseQuote, ""},
			{CloseParen, ""},
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
