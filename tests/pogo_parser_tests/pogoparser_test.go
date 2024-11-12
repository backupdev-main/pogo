package pogo_parser_tests

import (
	"fmt"
	"os"
	"pogo/src/lexer"
	"pogo/src/parser"
	"pogo/src/virtualmachine"
	"testing"
)

func TestParser(t *testing.T) {
	fmt.Println("Test Pogo Parser")
	inputFile := "simple.pogo"
	input, err := os.ReadFile(inputFile)

	if err != nil {
		t.Fatalf("Error reading input: %v", err)
	}

	lex := lexer.NewLexer(input)
	p := parser.NewParser(lex)

	err = p.ParseProgram()
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	} else {
		fmt.Println("Input successfully parsed!")
	}

	virtualMachine := virtualmachine.NewVirtualMachine(p.CodeGenerator.Quads, p.CodeGenerator.MemoryManager)
	if err := virtualMachine.Execute(); err != nil {
		fmt.Println("Error during execution", err)
	}
}

//func TestSimpleParser(t *testing.T) {
//	fmt.Println("Test Pogo Parser")
//	inputFile := "simple.pogo"
//	input, err := os.ReadFile(inputFile)
//
//	if err != nil {
//		t.Fatalf("Error reading input: %v", err)
//	}
//
//	lex := lexer.NewLexer(input)
//	p := parser.NewParser(lex)
//
//	err = p.ParseProgram()
//	if err != nil {
//		fmt.Println("Parse error:", err)
//	} else {
//		fmt.Println("Input successfully parsed!")
//	}
//}
