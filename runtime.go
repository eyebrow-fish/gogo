package gogo

import (
	"fmt"
	"strconv"
)

func Go(tree SyntaxTree) error {
	scopeData := make(map[string]interface{})

	for _, token := range tree.Children {
		switch token.Type {
		case BuiltinFunction:
			builtinFunctions[token.Data](token.Children)
		case VariableAssignment:
			if _, ok := scopeData[token.Children[0].Data]; !ok {
				assignVariable(scopeData, token)

				continue
			}

			return fmt.Errorf("variable already exists: %s", token.Children[0].Data)
		case VariableReassignment:
			if _, ok := scopeData[token.Children[0].Data]; ok {
				assignVariable(scopeData, token)

				continue
			}

			return fmt.Errorf("variable does not exist: %s", token.Children[0].Data)
		}
	}

	return nil
}

func assignVariable(scopeData map[string]interface{}, token SyntaxTree) {
	variableIdentifier := token.Children[0].Data
	value := token.Children[1]

	var variableData interface{}
	switch value.Type {
	case LiteralString:
		variableData = value.Data
	case LiteralInteger:
		variableData, _ = strconv.Atoi(value.Data)
	case LiteralBool:
		variableData, _ = strconv.ParseBool(value.Data)
	}

	scopeData[variableIdentifier] = variableData
}

var builtinFunctions = map[string]func([]SyntaxTree){
	"print": func(tree []SyntaxTree) {
		print(tree[0].Data)
	},
}
