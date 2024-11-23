package pogo_parser_tests

import (
	"fmt"
	"os"
	"pogo/src/lexer"
	"pogo/src/parser"
	"pogo/src/storer"
	"testing"
)

func TestParser(t *testing.T) {
	fmt.Println("Test Pogo Parser")
	inputFile := os.Args[len(os.Args)-1]
	//inputFile := "simple.pogo"
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
	}

	if err := storer.SaveCompiledData(p.CodeGenerator.Quads, p.SymbolTable, p.CodeGenerator.MemoryManager, "test.pbin"); err != nil {
		fmt.Println(err)
		return
	}

	vm, err := storer.LoadCompiledData("test.pbin")
	if err := vm.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

//func TestParserFibo(t *testing.T) {
//	fmt.Println("Test Pogo Parser")
//	inputFile := "fibo.pogo"
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
//		return
//	} else {
//		fmt.Println("Input successfully parsed!")
//	}
//
//	if err := storer.SaveCompiledData(p.CodeGenerator.Quads, p.SymbolTable, p.CodeGenerator.MemoryManager, "test.pbin"); err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	vm, err := storer.LoadCompiledData("test.pbin")
//	// fmt.Println(vm)
//	if err := vm.Execute(); err != nil {
//		fmt.Println(err)
//		return
//	}
//
//}
