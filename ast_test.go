package gogo

import (
	"reflect"
	"testing"
)

func TestBuildTrees(t *testing.T) {
	type args struct {
		tokens []Token
	}
	tests := []struct {
		name    string
		args    args
		want    *SyntaxTree
		wantErr bool
	}{
		{
			"variables",
			args{[]Token{{Identifier, "foo"}, {Type: Assignment}, {Literal, "123"}}},
			&SyntaxTree{Type: Program, Children: []SyntaxTree{{
				Type:     VariableAssignment,
				Children: []SyntaxTree{{Type: VariableIdentifier, Data: "foo"}, {Type: LiteralInteger, Data: "123"}},
			}}},
			false,
		},
		{
			"variable reassignment",
			args{[]Token{
				{Identifier, "foo"},
				{Type: Assignment},
				{Literal, "123"},
				{Type: Newline},
				{Identifier, "foo"},
				{Type: Reassignment},
				{Type: OpenQuote},
				{Literal, "bar"},
				{Type: CloseQuote},
				{Type: Semicolon},
				{Identifier, "foo"},
				{Type: Reassignment},
				{Literal, "true"},
			}},
			&SyntaxTree{Type: Program, Children: []SyntaxTree{
				{
					Type:     VariableAssignment,
					Children: []SyntaxTree{{Type: VariableIdentifier, Data: "foo"}, {Type: LiteralInteger, Data: "123"}},
				},
				{
					Type:     VariableReassignment,
					Children: []SyntaxTree{{Type: VariableIdentifier, Data: "foo"}, {Type: LiteralString, Data: "bar"}},
				},
				{
					Type:     VariableReassignment,
					Children: []SyntaxTree{{Type: VariableIdentifier, Data: "foo"}, {Type: LiteralBool, Data: "true"}},
				},
			}},
			false,
		},
		{
			"noob",
			args{[]Token{
				{Identifier, "print"},
				{Type: OpenParen},
				{Type: OpenQuote},
				{Literal, "Hello, World!"},
				{Type: CloseQuote},
				{Type: CloseParen},
			}},
			&SyntaxTree{Type: Program, Children: []SyntaxTree{{
				BuiltinFunction,
				"print",
				[]SyntaxTree{{Type: LiteralString, Data: "Hello, World!"}},
			}}},
			false,
		},
		{
			"multi-arg noob",
			args{[]Token{
				{Identifier, "print"},
				{Type: OpenParen},
				{Type: OpenQuote},
				{Literal, "Hello,"},
				{Type: CloseQuote},
				{Type: Comma},
				{Type: Whitespace},
				{Type: OpenQuote},
				{Literal, "World!"},
				{Type: CloseQuote},
				{Type: CloseParen},
			}},
			&SyntaxTree{Type: Program, Children: []SyntaxTree{{
				BuiltinFunction,
				"print",
				[]SyntaxTree{{Type: LiteralString, Data: "Hello,"}, {Type: LiteralString, Data: "World!"}},
			}}},
			false,
		},
		{
			"shell list dir",
			args{[]Token{{Bang, "ls -al"}}},
			&SyntaxTree{Type: Program, Children: []SyntaxTree{{
				Type: ShellCmd,
				Data: "ls -al",
			}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildTrees(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildTrees() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildTrees() got = %v, want %v", got, tt.want)
			}
		})
	}
}
