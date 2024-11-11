package semantic

import (
	"fmt"
	"pogo/src/shared"
)

// Variable represents a variable in our symbol table
type Variable struct {
	Name    string
	Type    shared.Type
	Line    int
	Column  int
	Address int
}

// Function represents a function in our function directory
type Function struct {
	Name       string
	Parameters []Variable
	Line       int
	Column     int
	StartQuad  int
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

func (st *SymbolTable) GetType(name string) (shared.Type, error) {
	if st.variables[st.currentScope] == nil {
		st.variables[st.currentScope] = make(map[string]interface{})
	}

	if symbol, exists := st.variables[st.currentScope][name]; exists {
		switch v := symbol.(type) {
		case Variable:
			return v.Type, nil
		default:
			return shared.TypeError, fmt.Errorf("symbol '%s' is not a variable", name)
		}
	}

	if st.currentScope != "global" {
		if symbol, exists := st.variables["global"][name]; exists {
			switch v := symbol.(type) {
			case Variable:
				return v.Type, nil
			default:
				return shared.TypeError, fmt.Errorf("symbol '%s' is not a variable", name)
			}
		}
	}

	return shared.TypeError, fmt.Errorf("variable '%s' not declared in accessible scope", name)
}

func (st *SymbolTable) GetVariableAddress(name string) (int, error) {
	// fmt.Println("This is the current scope: ", st.currentScope)

	value, exists := st.variables[st.currentScope][name]
	if exists {
		// then we should look for the value in global scope
		if v, ok := value.(Variable); ok {
			return v.Address, nil
		}
	}

	if st.currentScope != "global" {
		if value, exists := st.variables["global"][name]; exists {
			if v, ok := value.(Variable); ok {
				return v.Address, nil
			}
		}
	}

	return -1, fmt.Errorf("error retrieving address for '%v", name)
}

func (st *SymbolTable) AddVariable(name string, varType shared.Type, line, column int, addr int) error {
	// Don't allow declaring string variables
	if st.variables[st.currentScope] == nil {
		st.variables[st.currentScope] = make(map[string]interface{})
	}

	if varType == shared.TypeString {
		return fmt.Errorf("line %d: cannot declare string variables, strings are only allowed in print statements", line)
	}

	if _, exists := st.variables[st.currentScope][name]; exists {
		return fmt.Errorf("line %d: symbol '%s' already declared in current scope", line, name)
	}

	st.variables[st.currentScope][name] = Variable{
		Name:    name,
		Type:    varType,
		Line:    line,
		Column:  column,
		Address: addr,
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
		StartQuad:  -1,
	}

	// Add parameters to function scope
	for _, param := range params {
		st.variables[name][param.Name] = param
	}

	return nil
}

func (st *SymbolTable) UpdateFunctionStartQuad(functionName string, start int) error {
	function, ok := st.variables["global"][functionName].(Function)
	if !ok {
		return fmt.Errorf("function %s not found", functionName)
	}
	function.StartQuad = start

	st.variables["global"][functionName] = function
	return nil
}

func (st *SymbolTable) GetFunctionStartQuad(functionName string) (int, error) {
	function, ok := st.variables["global"][functionName].(Function)
	fmt.Println("Function Name", functionName, function)
	if !ok {
		return -1, fmt.Errorf("function %s not found", functionName)
	}

	return function.StartQuad, nil
}

// Functionality to be added later :/
//func (st *SymbolTable) UpdateFunctionMemoryRequirements(functionName string) error {
//	function, ok := st.variables["global"][functionName].(Function)
//	if !ok {
//		return fmt.Errorf("function %s not found", functionName)
//	}
//
//	// Count local variables in function scope
//	for _, symbol := range st.variables[functionName] {
//		if variable, ok := symbol.(Variable); ok {
//			switch variable.Type {
//			case shared.TypeInt:
//				function.IntVarsCount++
//			case shared.TypeFloat:
//				function.FloatVarsCount++
//			}
//		}
//	}
//
//	// Update function in global scope
//	st.variables["global"][functionName] = function
//	return nil
//}

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

	// If not in current scope, check global
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

func (st *SymbolTable) ValidateFunctionCall(funcName string, line int, args []shared.Type) error {
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
		if paramVar.Type == shared.TypeFloat && argType == shared.TypeInt {
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

func (st *SymbolTable) GetFunctionInfo(functionName string) (*Function, error) {
	symbol, exists := st.variables["global"][functionName]
	if !exists {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	funcInfo, ok := symbol.(Function)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", functionName)
	}

	return &funcInfo, nil
}

func (st *SymbolTable) GetScope() string {
	return st.currentScope
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
