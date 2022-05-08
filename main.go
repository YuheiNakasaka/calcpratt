package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

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
	for l.position < len(l.input) && isDigit(l.input[l.position]) {
		l.position++
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

const (
	_ int = iota
	LOWEST
	SUM     // +
	PRODUCT // *
	PREFIX  // -X or !X
)

var precedences = map[TokenType]int{
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
}

type Expression interface {
	TokenLiteral() string
	String() string
}

type IntegerLiteral struct {
	Token Token
	Value int64
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left, ie.Operator, ie.Right.String())
}

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
}

func New(l *Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() Expression {
	var expression Expression
	for p.currentToken.Type != EOF {
		expression = p.parseExpression(LOWEST)
		if expression != nil {
			return expression
		}
	}
	return expression
}

func (p *Parser) parseExpression(precedence int) Expression {
	var leftExp Expression

	switch p.currentToken.Type {
	case INT:
		lit := &IntegerLiteral{Token: p.currentToken}
		value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
		if err != nil {
			panic(err)
		}
		lit.Value = value
		leftExp = lit
	case MINUS:
		expression := &PrefixExpression{
			Token:    p.currentToken,
			Operator: p.currentToken.Literal,
		}
		p.nextToken()
		expression.Right = p.parseExpression(PREFIX)
		leftExp = expression
	}

	for p.peekToken.Literal != SEMICOLON && p.peekToken.Literal != EOF && precedence < p.peekPrecedence() {
		peekTokenType := p.peekToken.Type
		p.nextToken()
		switch peekTokenType {
		case PLUS:
			expression := &InfixExpression{
				Token:    p.currentToken,
				Operator: p.currentToken.Literal,
				Left:     leftExp,
			}
			p.nextToken()
			expression.Right = p.parseExpression(SUM)
			leftExp = expression
		case MINUS:
			expression := &InfixExpression{
				Token:    p.currentToken,
				Operator: p.currentToken.Literal,
				Left:     leftExp,
			}
			p.nextToken()
			expression.Right = p.parseExpression(SUM)
			leftExp = expression
		case ASTERISK:
			expression := &InfixExpression{
				Token:    p.currentToken,
				Operator: p.currentToken.Literal,
				Left:     leftExp,
			}
			p.nextToken()
			expression.Right = p.parseExpression(PRODUCT)
			leftExp = expression
		case SLASH:
			expression := &InfixExpression{
				Token:    p.currentToken,
				Operator: p.currentToken.Literal,
				Left:     leftExp,
			}
			p.nextToken()
			expression.Right = p.parseExpression(PRODUCT)
			leftExp = expression
		}
	}

	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := Lexer{input: line, position: 0}
		parser := New(&lexer)
		ast := parser.ParseProgram()
		fmt.Printf("%+v\n", ast)
	}
}

func main() {
	Start(os.Stdin, os.Stdout)
}
