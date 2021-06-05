package gogo

func Go(tree SyntaxTree) {
	for _, token := range tree.Children {
		switch token.Type {
		case BuiltinFunction:
			builtinFunctions[token.Data](token.Children)
		}
	}
}

var builtinFunctions = map[string]func([]SyntaxTree){
	"print": func(tree []SyntaxTree) {
		print(tree[0].Data)
	},
}
