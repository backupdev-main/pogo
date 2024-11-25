# Pogo Compiler
<p align="center">
  <img src="/pogologo.png" width="400" style="display: block; margin: 0 auto;" alt="logo">
</p>

A compiler for the **Pogo** programming language, implemented in Go. This hybrid compiler performs **lexical analysis**, **syntax parsing**, **semantic validation**, **intermediate code generation** and **code execution**.

## Features

- **Lexical Analysis**: Token recognition powered by GOCC
- **Recursive Descent Parsing**: Efficient parsing for context-free grammar
- **Symbol Table Management**: Tracks identifiers and scope for variables and functions
- **Type Checking**: Uses a semantic cube for enforcing type rules
- **Data Type Support**: Handles basic data types like `int` and `float`
- **Function Declarations & Calls**: Supports defining and invoking functions
- **Control Structures**: Implements control flow with `if` and `while` statements

## Examples

### Important Notes
- **Variable Declaration**: Variables can only be declared at the start of the program after program name and before function declarations, variables can also be declared within functions before any statement.
- **Functions**: Currently functions are only void functions.
- **Variable Types**: The program currently only handles ints and floats, booleans and comparisons are handled as ints.

### Example 1 Factorial
```
program Factorial;

var result : float;
var x : int;

begin
    // Comments work as well!!!!
    /*
        multiline comments also work!
        look!
    */
    result = 1;
    x = 5;
    while (x > 0) {
        result = result * x;
        x = x - 1;
    }
    print("This is the result", result)
end
```

### Example 2 Fibonacci

```
program recursiveFibo;

var result : int;

func fib(n : int) {
    var temp1, temp2 : int;

    if(n < 2) {
        result = n;
    } else {
        temp1 = n - 1;
        temp2 = n - 2;
        fib(temp1)
        temp1 = result;
        fib(temp2)
        temp2 = result;
        result = temp1 + temp2;
    }
};

begin
    fib(30)
    print("This is the result", result)
end
```

## How to Run
How to Run
The main.go script demonstrates the compilation and execution process:

### Lexical Analysis:

Create a lexer object by passing an input file
Generate tokens through lexical analysis


### Parsing:

Initialize the parser with the lexer tokens
Perform recursive descent parsing


### Compilation:

Generate a binary file containing compiled data
Serialize necessary information for VM execution


### Execution:

Load the compiled binary file
Execute the virtual machine


### Example Workflow

```
lex := lexer.NewLexer(inputFile)

// Initialize parser
parser := parser.NewParser(lex)

// Parse the program
parser.ParseProgram()

// Save compiled data
storer.SaveCompiledData(parser.CodeGenerator.Quads, parser.SymbolTable, parser.CodeGenerator.MemoryManager, "output.pbin")

// Load and execute
vm := storer.LoadCompiledData("output.pbin")
vm.Execute()
```
