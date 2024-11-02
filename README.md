# Pogo Compiler
<p align="center">
  <img src="/pogologo.png" width="400" style="display: block; margin: 0 auto;" alt="logo">
</p>

A compiler for the **Pogo** programming language, implemented in Go. This hybrid compiler performs **lexical analysis**, **syntax parsing**, **semantic validation**, "intermediate code generation" and "code execution".

## Features

- **Lexical Analysis**: Token recognition powered by GOCC
- **Recursive Descent Parsing**: Efficient parsing for context-free grammar
- **Symbol Table Management**: Tracks identifiers and scope for variables and functions
- **Type Checking**: Uses a semantic cube for enforcing type rules
- **Data Type Support**: Handles basic data types like `int` and `float`
- **Function Declarations & Calls**: Supports defining and invoking functions
- **Control Structures**: Implements control flow with `if` and `while` statements
