package tests

import (
	"fmt"
	"os"
	"pogo/src/lexer"
	"pogo/src/parser"
	"testing"
)

func TestParser(t *testing.T) {
	fmt.Println("Test1")
	inputFile := "patito.pogo"
	input, err := os.ReadFile(inputFile)

	if err != nil {
		t.Fatalf("Error reading input: %v", err)
	}

	lex := lexer.NewLexer(input)
	p := parser.NewParser()

	_, err = p.Parse(lex)
	if err != nil {
		fmt.Println("Parse error:", err)
	} else {
		fmt.Println("Input successfully parsed!")
	}
}

func TestParserSimple(t *testing.T) {
	inputFile := "simple_patito.pogo"
	input, err := os.ReadFile(inputFile)

	if err != nil {
		t.Fatalf("Error reading input: %v", err)
	}

	lex := lexer.NewLexer(input)
	p := parser.NewParser()

	_, err = p.Parse(lex)
	if err != nil {
		fmt.Println("Parse error:", err)
	} else {
		fmt.Println("Input successfully parsed!")
	}
}

func TestParserNested(t *testing.T) {
	inputFile := "nested_patito.pogo"
	input, err := os.ReadFile(inputFile)

	if err != nil {
		t.Fatalf("Error reading input: %v", err)
	}

	lex := lexer.NewLexer(input)
	p := parser.NewParser()

	_, err = p.Parse(lex)

	if err != nil {
		fmt.Println("Parse error:", err)
	} else {
		fmt.Println("Input successfully parsed!")
	}
}
