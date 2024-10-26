package pogo_parser_tests

import (
	"fmt"
	"os"
	"pogo/src/lexer"
	"pogo/src/pogo_parser"
	"testing"
)

func TestParser(t *testing.T) {
	fmt.Println("Test Pogo Parser")
	inputFile := "pogoPatito.pogo"
	input, err := os.ReadFile(inputFile)

	if err != nil {
		t.Fatalf("Error reading input: %v", err)
	}

	lex := lexer.NewLexer(input)
	p := pogo_parser.NewParser(lex)

	err = p.ParseProgram()
	if err != nil {
		fmt.Println("Parse error:", err)
	} else {
		fmt.Println("Input successfully parsed!")
	}
}
