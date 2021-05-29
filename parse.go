package gogo

type parseState struct {
	tokens       []Token
	visibleToken *Token
}

func (ps *parseState) appendVisibleToken(terminatingC rune) {
	var terminatingToken *Token
	if terminatingC == ' ' {
		terminatingToken = &Token{Whitespace, ""}
	} else if terminatingC == ';' {
		terminatingToken = &Token{Semicolon, ""}
	}

	if ps.visibleToken != nil {
		if terminatingToken != nil {
			ps.tokens = append(ps.tokens, *ps.visibleToken, *terminatingToken)
		} else {
			ps.tokens = append(ps.tokens, *ps.visibleToken)
		}

		ps.visibleToken = nil
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
		case ' ', ';':
			ps.appendVisibleToken(c)
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
			if ps.visibleToken != nil && ps.visibleToken.Type == Literal {
				ps.visibleToken.Data += string(c)

				if lastC {
					ps.appendVisibleToken(c)
				}

				continue
			}

			ps.visibleToken = &Token{Literal, string(c)}
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
	Literal TokenType = iota
	Comment
	Whitespace
	Assignment
	Reassignment
	Equals
	Plus
	Minus
	Newline
	Semicolon
	OpenParen
	CloseParen
	OpenBracket
	CloseBracket
)
