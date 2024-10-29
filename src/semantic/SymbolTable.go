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
	Variables map[string]Variable
	Functions map[string]Function
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Variables: make(map[string]Variable),
		Functions: make(map[string]Function),
	}
}

func (st *SymbolTable) GetType(name string) (Type, error) {
	variable, exists := st.Variables[name]

	if !exists {
		return TypeError, fmt.Errorf("variable '%s' not declared", name)
	}

	return variable.Type, nil
}

func (st *SymbolTable) AddVariable(name string, varType Type, line, column int) error {
	// Don't allow declaring string variables
	if varType == TypeString {
		return fmt.Errorf("line %d: cannot declare string variables, strings are only allowed in print statements", line)
	}

	if _, exists := st.Variables[name]; exists {
		return fmt.Errorf("line %d: variable '%s' already declared", line, name)
	}

	st.Variables[name] = Variable{
		Name:   name,
		Type:   varType,
		Line:   line,
		Column: column,
	}

	return nil
}

func (st *SymbolTable) AddFunction(name string, params []Variable, line, column int) error {
	if _, isVar := st.Variables[name]; isVar {
		return fmt.Errorf("line %d: cannot declare function '%s', name already used by a variable", line, name)
	}

	if _, exists := st.Functions[name]; exists {
		return fmt.Errorf("line %d: function '%s' already declared", line, name)
	}

	st.Functions[name] = Function{
		Name:       name,
		Parameters: params,
		Line:       line,
		Column:     column,
	}

	return nil
}

func (st *SymbolTable) ValidatePrintStatement(items []Type) error {
	for i, itemType := range items {
		if itemType != TypeInt && itemType != TypeFloat && itemType != TypeString {
			return fmt.Errorf("invalid type in print statement at position %d: only int, float, and string literals are allowed", i+1)
		}
	}
	return nil
}

func (st *SymbolTable) ValidateVarAssignment(varName string, line int) error {
	// First check if variable exists
	_, exists := st.Variables[varName]
	if !exists {
		return fmt.Errorf("line %d: undefined variable '%s'", line, varName)
	}

	return nil
}

func (st *SymbolTable) ValidateFunctionCall(funcName string, line int) error {
	_, exists := st.Functions[funcName]
	if !exists {
		return fmt.Errorf("line %d: undefined function '%s'", line, funcName)
	}

	return nil
}
