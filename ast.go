package gogo

import (
	"errors"
	"fmt"
	"strconv"
)

// Effectively a set for fast indexing
var builtinFunctionNames = map[string]interface{}{
	"print": nil,
}

type treeState struct {
	currentIdentifier string
	currentTree       *SyntaxTree
	inStringLiteral   bool
	inFunctionParen   bool
}

func BuildTrees(tokens []Token) (*SyntaxTree, error) {
	var tree SyntaxTree
	var ts treeState

	for _, token := range tokens {
		var lastRootTree *SyntaxTree
		if len(tree.Children) > 0 {
			lastRootTree = &tree.Children[len(tree.Children)-1]
		}

		switch token.Type {
		case Identifier:
			if ts.currentIdentifier != "" {
				return nil, fmt.Errorf("unexpected identifier: %s", ts.currentIdentifier)
			}

			ts.currentIdentifier = token.Data
		case Assignment, Reassignment:
			if ts.currentIdentifier == "" {
				return nil, errors.New(`expected identifier, got ":="`)
			}

			var assignmentType TreeType
			if token.Type == Assignment {
				assignmentType = VariableAssignment
			} else {
				assignmentType = VariableReassignment
			}

			tree.Children = append(tree.Children, SyntaxTree{
				Type:     assignmentType,
				Children: []SyntaxTree{{Type: VariableIdentifier, Data: ts.currentIdentifier}},
			})

			ts.currentIdentifier = ""
		case Literal:
			// Literal typing
			var literalType TreeType
			if ts.inStringLiteral {
				literalType = LiteralString
			} else {
				if token.Data == "true" || token.Data == "false" {
					literalType = LiteralBool
				} else if _, err := strconv.Atoi(token.Data); err == nil {
					literalType = LiteralInteger
				} else {
					return nil, fmt.Errorf("could not literal %s: %v", token.Data, err)
				}
			}

			// Function call
			if ts.currentTree != nil && ts.currentTree.Type == BuiltinFunction {
				ts.currentTree.Children = append(ts.currentTree.Children, SyntaxTree{Type: literalType, Data: token.Data})
			}

			// Everything below require a lastRootTree
			if lastRootTree == nil {
				continue
			}

			// Variable assignment
			deepestChild := lastRootTree.Children[len(lastRootTree.Children)-1]
			if lastRootTree.Type != VariableAssignment && lastRootTree.Type != VariableReassignment {
				return nil, fmt.Errorf("unexpected literal: %s", token.Data)
			}
			if deepestChild.Type == LiteralInteger || deepestChild.Type == LiteralString {
				return nil, fmt.Errorf(`unexpected literal "%s" after literal "%s"`, token.Data, deepestChild.Data)
			}

			lastRootTree.Children = append(lastRootTree.Children, SyntaxTree{Type: literalType, Data: token.Data})
			ts.currentIdentifier = ""
		case OpenParen:
			if ts.inFunctionParen {
				return nil, fmt.Errorf("unexpected (, already in function parenthesis")
			}

			if _, ok := builtinFunctionNames[ts.currentIdentifier]; !ok {
				return nil, fmt.Errorf("unknown builtin function: %s", ts.currentIdentifier)
			}

			ts.currentTree = &SyntaxTree{Type: BuiltinFunction, Data: ts.currentIdentifier}

			ts.inFunctionParen = true
		case CloseParen:
			if !ts.inFunctionParen || ts.currentTree == nil {
				return nil, fmt.Errorf("unexpected ), no matching open parenthesis")
			}

			tree.Children = append(tree.Children, *ts.currentTree)

			ts.inFunctionParen = false
			ts.currentTree = nil
		case Bang:
			tree.Children = append(tree.Children, SyntaxTree{Type: ShellCmd, Data: token.Data})
		case Comma:
		case OpenQuote:
			ts.inStringLiteral = true
		case CloseQuote:
			ts.inStringLiteral = false
		case Newline, Semicolon:
			ts.currentIdentifier = ""
		default: // Just ignore lol
		}
	}

	return &tree, nil
}

type SyntaxTree struct {
	Type     TreeType
	Data     string
	Children []SyntaxTree
}

type TreeType uint8

const (
	Program TreeType = iota
	VariableAssignment
	VariableReassignment
	VariableIdentifier
	LiteralInteger
	LiteralString
	LiteralBool
	BuiltinFunction
	ShellCmd
)
