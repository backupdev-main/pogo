package semantic

import (
	"fmt"
)

type Type int

const (
	TypeInt Type = iota
	TypeFloat
	TypeString // Added string type
	TypeError
)

func (t Type) String() string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeString:
		return "string"
	default:
		return "error"
	}
}

// Variable represents a variable in our symbol table
type Variable struct {
	Name   string
	Type   Type
	Line   int
	Column int
}

// Function represents a function in our function directory
type Function struct {
	Name       string
	Parameters []Variable
	Line       int
	Column     int
}

type SymbolTable struct {
	variables    map[string]map[string]interface{}
	scopeStack   []string
	currentScope string
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		variables: make(map[string]map[string]interface{}),
	}

	st.variables["global"] = make(map[string]interface{})
	st.scopeStack = append(st.scopeStack, "global")
	st.currentScope = st.scopeStack[len(st.scopeStack)-1]
	return st
}

func (st *SymbolTable) GetType(name string) (Type, error) {
	if st.variables[st.currentScope] == nil {
		st.variables[st.currentScope] = make(map[string]interface{})
	}

	if symbol, exists := st.variables[st.currentScope][name]; exists {
		switch v := symbol.(type) {
		case Variable:
			return v.Type, nil
		default:
			return TypeError, fmt.Errorf("symbol '%s' is not a variable", name)
		}
	}

	if st.currentScope != "global" {
		if symbol, exists := st.variables["global"][name]; exists {
			switch v := symbol.(type) {
			case Variable:
				return v.Type, nil
			default:
				return TypeError, fmt.Errorf("symbol '%s' is not a variable", name)
			}
		}
	}

	return TypeError, fmt.Errorf("variable '%s' not declared in accessible scope", name)
}

func (st *SymbolTable) AddVariable(name string, varType Type, line, column int) error {
	// Don't allow declaring string variables
	if st.variables[st.currentScope] == nil {
		st.variables[st.currentScope] = make(map[string]interface{})
	}

	if varType == TypeString {
		return fmt.Errorf("line %d: cannot declare string variables, strings are only allowed in print statements", line)
	}

	if _, exists := st.variables[st.currentScope][name]; exists {
		return fmt.Errorf("line %d: symbol '%s' already declared in current scope", line, name)
	}

	st.variables[st.currentScope][name] = Variable{
		Name:   name,
		Type:   varType,
		Line:   line,
		Column: column,
	}

	return nil
}

func (st *SymbolTable) AddFunction(name string, params []Variable, line, column int) error {
	if _, exists := st.variables["global"][name]; exists {
		return fmt.Errorf("line %d: symbol '%s' already declared", line, name)
	}

	st.variables[name] = make(map[string]interface{})

	// Add function to global scope
	st.variables["global"][name] = Function{
		Name:       name,
		Parameters: params,
		Line:       line,
		Column:     column,
	}

	// Add parameters to function scope
	for _, param := range params {
		st.variables[name][param.Name] = param
	}

	return nil
}

func (st *SymbolTable) ValidateVarAssignment(varName string, line int) error {
	// First check if variable exists
	if st.variables[st.currentScope] == nil {
		st.variables[st.currentScope] = make(map[string]interface{})
	}

	if symbol, exists := st.variables[st.currentScope][varName]; exists {
		if _, ok := symbol.(Variable); !ok {
			return fmt.Errorf("line %d: '%s' is not a variable", line, varName)
		}
		return nil
	}

	// If in function scope, check global
	if st.currentScope != "global" {
		if symbol, exists := st.variables["global"][varName]; exists {
			if _, ok := symbol.(Variable); !ok {
				return fmt.Errorf("line %d: '%s' is not a variable", line, varName)
			}
			return nil
		}
	}

	return fmt.Errorf("line %d: undefined variable '%s'", line, varName)
}

func (st *SymbolTable) ValidateFunctionCall(funcName string, line int, args []Type) error {
	symbol, exists := st.variables["global"][funcName]
	if !exists {
		return fmt.Errorf("line %d: undefined function '%s'", line, funcName)
	}

	function, ok := symbol.(Function)
	if !ok {
		return fmt.Errorf("line %d: '%s' is not a function", line, funcName)
	}

	if len(function.Parameters) != len(args) {
		return fmt.Errorf("line %d: function '%s' expects %d arguments but got %d",
			line, funcName, len(function.Parameters), len(args))
	}

	for i, paramVar := range function.Parameters {
		argType := args[i]

		// Direct type match
		if paramVar.Type == argType {
			continue
		}

		// Allow int -> float conversion
		if paramVar.Type == TypeFloat && argType == TypeInt {
			continue
		}

		return fmt.Errorf("line %d: invalid argument type for parameter '%s' in function '%s': expected %s, got %s",
			line, paramVar.Name, funcName, paramVar.Type, argType)
	}

	return nil
}

func (st *SymbolTable) ExitFunctionScope() {
	st.scopeStack = st.scopeStack[:len(st.scopeStack)-1]
	st.currentScope = "global"
}

func (st *SymbolTable) EnterFunctionScope(name string) error {
	if _, exists := st.variables["global"][name]; !exists {
		return fmt.Errorf("cannot enter scope of undefined function '%s'", name)
	}

	if _, ok := st.variables["global"][name].(Function); !ok {
		return fmt.Errorf("'%s' is not a function", name)
	}

	// Ensure the function scope map is initialized
	if st.variables[name] == nil {
		st.variables[name] = make(map[string]interface{})
	}

	st.scopeStack = append(st.scopeStack, name)
	st.currentScope = name
	return nil
}

func (st *SymbolTable) PrettyPrint() {
	fmt.Println("\n=== Symbol Table ===")

	// Print global scope first
	fmt.Println("\nGlobal Scope:")
	fmt.Println("-------------")
	for name, symbol := range st.variables["global"] {
		switch v := symbol.(type) {
		case Variable:
			fmt.Printf("Variable: %s\n", name)
			fmt.Printf("  Type: %s\n", v.Type)
			fmt.Printf("  Line: %d, Column: %d\n", v.Line, v.Column)
		case Function:
			fmt.Printf("Function: %s\n", name)
			fmt.Printf("  Parameters:\n")
			if len(v.Parameters) == 0 {
				fmt.Printf("    None\n")
			}
			for _, param := range v.Parameters {
				fmt.Printf("    - %s: %s\n", param.Name, param.Type)
			}
			fmt.Printf("  Line: %d, Column: %d\n", v.Line, v.Column)
		}
		fmt.Println()
	}

	// Print other scopes (function scopes)
	for scope, symbols := range st.variables {
		if scope != "global" {
			fmt.Printf("\nFunction Scope: %s\n", scope)
			fmt.Println("------------------")
			for name, symbol := range symbols {
				switch v := symbol.(type) {
				case Variable:
					fmt.Printf("Variable: %s\n", name)
					fmt.Printf("  Type: %s\n", v.Type)
					fmt.Printf("  Line: %d, Column: %d\n", v.Line, v.Column)
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("=== End Symbol Table ===")
}
