package main

import "testing"

func TestNextToken(t *testing.T) {
	input := `-1 + 2 * 3`
	lexer := Lexer{input: input, position: 0}

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{MINUS, "-"},
		{INT, "1"},
		{PLUS, "+"},
		{INT, "2"},
		{ASTERISK, "*"},
		{INT, "3"},
	}

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestParseProgram(t *testing.T) {
	inputs := []struct {
		input    string
		expected string
	}{
		{"1 + 2", "(1 + 2)"},
		{"-1 + 2 * 3", "((-1) + (2 * 3))"},
		{"1 - -2", "(1 - (-2))"},
		{"1 - -2 / 2 * 3 + 5", "((1 - (((-2) / 2) * 3)) + 5)"},
	}
	for _, tt := range inputs {
		lexer := Lexer{input: tt.input, position: 0}
		parser := New(&lexer)
		ast := parser.ParseProgram()
		if ast.String() != tt.expected {
			t.Fatalf("expected=%q, got=%q", tt.expected, ast.String())
		}
	}
}
