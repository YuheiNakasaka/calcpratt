package main

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF       = "EOF"
	INT       = "INT"
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	SEMICOLON = ";"
)

type Lexer struct {
	input    string
	position int
}

func (l *Lexer) NextToken() Token {
	var token Token

	if l.position >= len(l.input) {
		return Token{Type: EOF, Literal: ""}
	}

	for l.position < len(l.input) {
		if l.input[l.position] == ' ' {
			l.position++
			continue
		}

		switch l.input[l.position] {
		case '+':
			token = Token{Type: PLUS, Literal: string(l.input[l.position])}
			l.position++
			return token
		case '-':
			token = Token{Type: MINUS, Literal: string(l.input[l.position])}
			l.position++
			return token
		case '*':
			token = Token{Type: ASTERISK, Literal: string(l.input[l.position])}
			l.position++
			return token
		case '/':
			token = Token{Type: SLASH, Literal: string(l.input[l.position])}
			l.position++
			return token
		case ';':
			token = Token{Type: SEMICOLON, Literal: string(l.input[l.position])}
			l.position++
			return token
		default:
			if isDigit(l.input[l.position]) {
				token = Token{Type: INT, Literal: l.readNumber()}
				return token
			}
			l.position++
		}
	}
	return token
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.input[l.position]) {
		l.position++
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func main() {
	input := "1 + 2 * 3;"
	lexer := Lexer{input: input, position: 0}
	fmt.Println(input)
	fmt.Println(lexer.NextToken())
	fmt.Println(lexer.NextToken())
	fmt.Println(lexer.NextToken())
	fmt.Println(lexer.NextToken())
	fmt.Println(lexer.NextToken())
	fmt.Println(lexer.NextToken())
	fmt.Println(lexer.NextToken())
}
