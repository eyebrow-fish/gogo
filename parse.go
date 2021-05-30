package gogo

import "strconv"

type parseState struct {
	tokens          []Token
	visibleToken    *Token
	inComment       bool
	inStringLiteral bool
}

func (ps *parseState) appendVisibleToken(terminatingC rune) {
	var terminatingToken *Token
	if terminatingC == ' ' {
		terminatingToken = &Token{Whitespace, ""}
	} else if terminatingC == ';' {
		terminatingToken = &Token{Semicolon, ""}
	} else if terminatingC == '\n' {
		terminatingToken = &Token{Newline, ""}
	}

	if ps.visibleToken != nil {
		if ps.visibleToken.Type == Identifier {
			if _, err := strconv.Atoi(ps.visibleToken.Data); err == nil {
				ps.visibleToken.Type = Literal
			} else if _, err := strconv.ParseBool(ps.visibleToken.Data); err == nil {
				ps.visibleToken.Type = Literal
			}
		}

		if terminatingToken != nil {
			ps.tokens = append(ps.tokens, *ps.visibleToken, *terminatingToken)
		} else {
			ps.tokens = append(ps.tokens, *ps.visibleToken)
		}

		ps.visibleToken = nil
	} else if terminatingToken != nil {
		ps.tokens = append(ps.tokens, *terminatingToken)
	}
}

func Parse(program string) []Token {
	if program == "" {
		return []Token{}
	}

	ps := parseState{}

	for i, c := range program {
		lastC := i == len(program)-1

		switch c {
		case '#':
			ps.inComment = true
		case '\n':
			ps.inComment = false
			ps.appendVisibleToken(c)
			continue
		}

		if ps.inComment {
			continue
		}

		if ps.inStringLiteral && c != '"' {
			ps.visibleToken.Data += string(c)
			continue
		}

		switch c {
		case ' ', ';':
			ps.appendVisibleToken(c)
		case '"':
			if ps.visibleToken != nil && ps.inStringLiteral {
				ps.inStringLiteral = false
				ps.appendVisibleToken(c)
				ps.tokens = append(ps.tokens, Token{CloseQuote, ""})
				continue
			}

			ps.tokens = append(ps.tokens, Token{OpenQuote, ""})
			ps.inStringLiteral = true
			ps.visibleToken = &Token{Literal, ""}
		case '(':
			ps.appendVisibleToken(c)
			ps.tokens = append(ps.tokens, Token{OpenParen, ""})
		case ')':
			ps.appendVisibleToken(c)
			ps.tokens = append(ps.tokens, Token{CloseParen, ""})
		case ':':
			if ps.visibleToken == nil {
				ps.visibleToken = &Token{Assignment, ""}
			}
		case '=':
			if ps.visibleToken != nil {
				if ps.visibleToken.Type == Assignment {
					continue
				}

				if ps.visibleToken.Type == Reassignment {
					ps.visibleToken = &Token{Equals, ""}
				}
			}

			ps.visibleToken = &Token{Reassignment, ""}
		default:
			if ps.visibleToken != nil && ps.visibleToken.Type == Identifier {
				ps.visibleToken.Data += string(c)

				if lastC {
					ps.appendVisibleToken(c)
				}

				continue
			}

			ps.visibleToken = &Token{Identifier, string(c)}
		}
	}

	return ps.tokens
}

type Token struct {
	Type TokenType
	Data string
}

type TokenType uint8

const (
	Identifier TokenType = iota
	Literal
	Whitespace
	Assignment
	Reassignment
	Equals
	Newline
	Semicolon
	OpenQuote
	CloseQuote
	OpenParen
	CloseParen
)
