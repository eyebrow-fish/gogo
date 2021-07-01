package gogo

import (
	"fmt"
	"strconv"
	"strings"
)

func Go(tree SyntaxTree) error {
	scopeData := make(scopeData)

	for _, token := range tree.Children {
		switch token.Type {
		case BuiltinFunction:
			builtinFunctions[token.Data](scopeData, token.Children)
		case VariableAssignment:
			if _, ok := scopeData[token.Data]; !ok {
				assignVariable(scopeData, token)

				continue
			}

			return fmt.Errorf("variable already exists: %s", token.Children[0].Data)
		case VariableReassignment:
			if _, ok := scopeData[token.Data]; ok {
				assignVariable(scopeData, token)

				continue
			}

			return fmt.Errorf("variable does not exist: %s", token.Children[0].Data)
		}
	}

	return nil
}

func assignVariable(scopeData map[string]interface{}, token SyntaxTree) {
	variableIdentifier := token.Data
	value := token.Children[0]

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

var builtinFunctions = map[string]func(scopeData, []SyntaxTree){
	"print": func(scopeData scopeData, tree []SyntaxTree) {
		var datas []string
		for _, t := range tree {
			var data string
			if t.Type == VariableIdentifier {
				data = scopeData.getString(t.Data)
			} else {
				data = t.Data
			}

			datas = append(datas, data)
		}
		print(strings.Join(datas, " "))
	},
}

type scopeData map[string]interface{}

func (s scopeData) get(indent string) interface{} {
	return s[indent]
}

func (s scopeData) getString(indent string) string {
	data := s[indent]
	switch data.(type) {
	case string:
		return data.(string)
	case int:
		return strconv.Itoa(data.(int))
	case bool:
		return strconv.FormatBool(data.(bool))
	default:
		return "[UNKNOWN]"
	}
}
