package tests

import (
	"pogo/src/lexer"
	"pogo/src/token"
	"testing"
)

func TestTokenizer(t *testing.T) {
	type Test struct {
		expectedType    token.Type
		expectedLiteral string
	}

	const input = `
		program pablillo;
		var testVar, testVar1 : int;
		var test_var : float;
		
		begin
			testVar = 5;
			testVar1 = 7;
			if (testVar > testVar1) {
				print(testVar);
			}

			while (testVar == testVar1) {
				testVar = testVar - 1;
			}
		end
    `

	tests := []Test{
		{token.TokMap.Type("kwdProgram"), "program"},
		{token.TokMap.Type("id"), "pablillo"},
		{token.TokMap.Type("terminator"), ";"},
		{token.TokMap.Type("kwdVars"), "var"},
		{token.TokMap.Type("id"), "testVar"},
		{token.TokMap.Type("repeatTerminator"), ","},
		{token.TokMap.Type("id"), "testVar1"},
		{token.TokMap.Type("typeAssignOp"), ":"},
		{token.TokMap.Type("type"), "int"},
		{token.TokMap.Type("terminator"), ";"},
		{token.TokMap.Type("kwdVars"), "var"},
		{token.TokMap.Type("id"), "test_var"},
		{token.TokMap.Type("typeAssignOp"), ":"},
		{token.TokMap.Type("type"), "float"},
		{token.TokMap.Type("terminator"), ";"},
		{token.TokMap.Type("kwdBegin"), "begin"},
		{token.TokMap.Type("id"), "testVar"},
		{token.TokMap.Type("assignOp"), "="},
		{token.TokMap.Type("intLit"), "5"},
		{token.TokMap.Type("terminator"), ";"},
		{token.TokMap.Type("id"), "testVar1"},
		{token.TokMap.Type("assignOp"), "="},
		{token.TokMap.Type("intLit"), "7"},
		{token.TokMap.Type("terminator"), ";"},
		{token.TokMap.Type("kwdIf"), "if"},
		{token.TokMap.Type("openParan"), "("},
		{token.TokMap.Type("id"), "testVar"},
		{token.TokMap.Type("relOp"), ">"},
		{token.TokMap.Type("id"), "testVar1"},
		{token.TokMap.Type("closeParan"), ")"},
		{token.TokMap.Type("openBrace"), "{"},
		{token.TokMap.Type("kwdPrint"), "print"},
		{token.TokMap.Type("openParan"), "("},
		{token.TokMap.Type("id"), "testVar"},
		{token.TokMap.Type("closeParan"), ")"},
		{token.TokMap.Type("terminator"), ";"},
		{token.TokMap.Type("closeBrace"), "}"},
		{token.TokMap.Type("kwdWhile"), "while"},
		{token.TokMap.Type("openParan"), "("},
		{token.TokMap.Type("id"), "testVar"},
		{token.TokMap.Type("relOp"), "=="},
		{token.TokMap.Type("id"), "testVar1"},
		{token.TokMap.Type("closeParan"), ")"},
		{token.TokMap.Type("openBrace"), "{"},
		{token.TokMap.Type("id"), "testVar"},
		{token.TokMap.Type("assignOp"), "="},
		{token.TokMap.Type("id"), "testVar"},
		{token.TokMap.Type("expressionOp"), "-"},
		{token.TokMap.Type("intLit"), "1"},
		{token.TokMap.Type("terminator"), ";"},
		{token.TokMap.Type("closeBrace"), "}"},
		{token.TokMap.Type("kwdEnd"), "end"},
	}

	l := lexer.NewLexer([]byte(input))

	for i, tt := range tests {
		tok := l.Scan()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected='%s', got='%s' at line %d, column %d", i, token.TokMap.Id(tt.expectedType), token.TokMap.Id(tok.Type), tok.Pos.Line, tok.Pos.Column)
		}

		if string(tok.Lit) != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected='%q', got='%q' at line %d, column %d",
				i, tt.expectedLiteral, string(tok.Lit), tok.Pos.Line, tok.Pos.Column)
		}
	}
}
