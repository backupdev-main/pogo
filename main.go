package main

import (
	"fmt"
	"os"
	"pogo/src/lexer"
	"pogo/src/parser"
	"pogo/src/storer"
)

func main() {
	inputFile := os.Args[len(os.Args)-1]
	input, err := os.ReadFile(inputFile)

	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
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
