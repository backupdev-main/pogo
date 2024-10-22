package main

import (
	"fmt"
	"os"
	"pogo/src/lexer"
	"pogo/src/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input-file>")
		return
	}
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
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
