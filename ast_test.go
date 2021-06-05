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
				Children: []SyntaxTree{{Type: VariableIdentifier, Data: "foo"}, {Type: LiteralValue, Data: "123"}},
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