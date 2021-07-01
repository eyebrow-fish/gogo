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

	appendCurrentTree := func() {
		tree.Children = append(tree.Children, *ts.currentTree)
		ts.currentTree = nil
	}

	for _, token := range tokens {
		switch token.Type {
		case Identifier:
			if ts.currentTree != nil {
				ts.currentTree.Children = append(ts.currentTree.Children, SyntaxTree{Type: VariableIdentifier, Data: token.Data})
				continue
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

			ts.currentTree = &SyntaxTree{Type: assignmentType, Data: ts.currentIdentifier}
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

			if ts.currentTree != nil {
				if ts.currentTree.Type == BuiltinFunction {
					ts.currentTree.Children = append(ts.currentTree.Children, SyntaxTree{Type: literalType, Data: token.Data})
				}
				if ts.currentTree.Type == VariableReassignment || ts.currentTree.Type == VariableAssignment {
					ts.currentTree.Children = append(ts.currentTree.Children, SyntaxTree{Type: literalType, Data: token.Data})
					appendCurrentTree()
				}
			}

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

			appendCurrentTree()

			ts.inFunctionParen = false
			ts.currentIdentifier = ""
		case Bang:
			tree.Children = append(tree.Children, SyntaxTree{Type: ShellCmd, Data: token.Data})
		case Comma:
		case OpenQuote:
			ts.inStringLiteral = true
		case CloseQuote:
			ts.inStringLiteral = false
		case Newline, Semicolon:
			if ts.currentTree != nil {
				appendCurrentTree()
			}

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
