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
			{OpenParen, ""},
			{Literal, "bar"},
			{CloseParen, ""},
		}},
		{"variable reassignment", args{`foo := "bar"; foo = "bux"`}, []Token{
			{Literal, "foo"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
			{OpenParen, ""},
			{Literal, "bar"},
			{CloseParen, ""},
			{Semicolon, ""},
			{Whitespace, ""},
			{Literal, "foo"},
			{Whitespace, ""},
			{Reassignment, ""},
			{Whitespace, ""},
			{OpenParen, ""},
			{Literal, "bux"},
			{CloseParen, ""},
		}},
		{"comments", args{"# hello!\nabc := 123"}, []Token{
			{Newline, ""},
			{Literal, "abc"},
			{Whitespace, ""},
			{Assignment, ""},
			{Whitespace, ""},
			{Literal, "123"},
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
