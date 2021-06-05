package gogo

import (
	"errors"
	"fmt"
	"strconv"
)

type treeState struct {
	currentIdentifier string
	inStringLiteral   bool
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
			if lastRootTree == nil {
				continue
			}

			deepestChild := lastRootTree.Children[len(lastRootTree.Children)-1]
			if lastRootTree.Type != VariableAssignment && lastRootTree.Type != VariableReassignment {
				return nil, fmt.Errorf("unexpected literal: %s", token.Data)
			}
			if deepestChild.Type == LiteralInteger || deepestChild.Type == LiteralString {
				return nil, fmt.Errorf(`unexpected literal "%s" after literal "%s"`, token.Data, deepestChild.Data)
			}

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

			lastRootTree.Children = append(lastRootTree.Children, SyntaxTree{Type: literalType, Data: token.Data})
			ts.currentIdentifier = ""
		case OpenParen:
			ts.inStringLiteral = true
		case CloseParen:
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
)
