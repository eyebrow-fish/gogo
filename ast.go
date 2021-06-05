package gogo

import (
	"errors"
	"fmt"
)

type treeState struct {
	currentIdentifier string
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
		case Assignment:
			if ts.currentIdentifier == "" {
				return nil, errors.New(`expected identifier, got ":="`)
			}

			tree.Children = append(tree.Children, SyntaxTree{
				Type:     VariableAssignment,
				Children: []SyntaxTree{{Type: VariableIdentifier, Data: ts.currentIdentifier}},
			})

			ts.currentIdentifier = ""
		case Literal:
			if lastRootTree == nil {
				continue
			}

			if lastRootTree.Type != VariableAssignment {
				return nil, fmt.Errorf("unexpected literal: %s", token.Data)
			}

			lastRootTree.Children = append(lastRootTree.Children, SyntaxTree{Type: LiteralValue, Data: token.Data})
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
	VariableIdentifier
	LiteralValue
)
