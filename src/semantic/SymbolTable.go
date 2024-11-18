package semantic

import (
	"fmt"
	"pogo/src/shared"
)

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
		case shared.Variable:
			return v.Type, nil
		default:
			return shared.TypeError, fmt.Errorf("symbol '%s' is not a variable", name)
		}
	}

	if st.currentScope != "global" {
		if symbol, exists := st.variables["global"][name]; exists {
			switch v := symbol.(type) {
			case shared.Variable:
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
		if v, ok := value.(shared.Variable); ok {
			return v.Address, nil
		}
	}

	if st.currentScope != "global" {
		if value, exists := st.variables["global"][name]; exists {
			if v, ok := value.(shared.Variable); ok {
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

	st.variables[st.currentScope][name] = shared.Variable{
		Name:    name,
		Type:    varType,
		Line:    line,
		Column:  column,
		Address: addr,
	}

	return nil
}

func (st *SymbolTable) AddFunction(name string, params []shared.Variable, line, column int) error {
	if _, exists := st.variables["global"][name]; exists {
		return fmt.Errorf("line %d: symbol '%s' already declared", line, name)
	}

	st.variables[name] = make(map[string]interface{})

	// Add function to global scope

	intCount := 0
	floatCount := 0

	// Count parameters by type
	for _, param := range params {
		switch param.Type {
		case shared.TypeInt:
			intCount++
		case shared.TypeFloat:
			floatCount++
		}
	}

	st.variables["global"][name] = shared.Function{
		Name:             name,
		Parameters:       params,
		Line:             line,
		Column:           column,
		StartQuad:        -1,
		IntVarsCounter:   intCount,
		FloatVarsCounter: floatCount,
	}

	// Add parameters to function scope
	for _, param := range params {
		st.variables[name][param.Name] = param
	}

	return nil
}

func (st *SymbolTable) IncrementFunctionVarCount(varType shared.Type) error {
	if st.currentScope == "global" {
		return fmt.Errorf("cannot increment function variable count in global scope")
	}

	function, ok := st.variables["global"][st.currentScope].(shared.Function)
	if !ok {
		return fmt.Errorf("current scope is not a function")
	}

	switch varType {
	case shared.TypeInt:
		function.IntVarsCounter++
	case shared.TypeFloat:
		function.FloatVarsCounter++
	default:
		return fmt.Errorf("unsupported variable type for counting")
	}

	// Update the function in global scope
	st.variables["global"][st.currentScope] = function
	return nil
}

func (st *SymbolTable) GetFunctionVarCounts(functionName string) (int, int, error) {
	function, ok := st.variables["global"][functionName].(shared.Function)
	if !ok {
		return 0, 0, fmt.Errorf("function %s not found", functionName)
	}

	return function.IntVarsCounter, function.FloatVarsCounter, nil
}

func (st *SymbolTable) UpdateFunctionStartQuad(functionName string, start int) error {
	function, ok := st.variables["global"][functionName].(shared.Function)
	if !ok {
		return fmt.Errorf("function %s not found", functionName)
	}
	function.StartQuad = start

	st.variables["global"][functionName] = function
	return nil
}

func (st *SymbolTable) GetFunctionStartQuad(functionName string) (int, error) {
	function, ok := st.variables["global"][functionName].(shared.Function)
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
	if st.variables[st.currentScope] == nil {
		st.variables[st.currentScope] = make(map[string]interface{})
	}

	if symbol, exists := st.variables[st.currentScope][varName]; exists {
		if _, ok := symbol.(shared.Variable); !ok {
			return fmt.Errorf("line %d: '%s' is not a variable", line, varName)
		}
		return nil
	}

	// If not in current scope, check global
	if st.currentScope != "global" {
		if symbol, exists := st.variables["global"][varName]; exists {
			if _, ok := symbol.(shared.Variable); !ok {
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

	function, ok := symbol.(shared.Function)
	if !ok {
		return fmt.Errorf("line %d: '%s' is not a function", line, funcName)
	}

	if len(function.Parameters) != len(args) {
		fmt.Println(function.Parameters)
		fmt.Println(args)
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

	if _, ok := st.variables["global"][name].(shared.Function); !ok {
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

func (st *SymbolTable) GetFunctionInfo(functionName string) (*shared.Function, error) {
	symbol, exists := st.variables["global"][functionName]
	if !exists {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	funcInfo, ok := symbol.(shared.Function)
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

	fmt.Println("\nGlobal Scope:")
	fmt.Println("-------------")
	for name, symbol := range st.variables["global"] {
		switch v := symbol.(type) {
		case shared.Variable:
			fmt.Printf("Variable: %s\n", name)
			fmt.Printf("  Type: %s\n", v.Type)
			fmt.Printf("  Line: %d, Column: %d\n", v.Line, v.Column)
		case shared.Function:
			fmt.Printf("Function: %s\n", name)
			fmt.Printf("  Parameters:\n")
			if len(v.Parameters) == 0 {
				fmt.Printf("    None\n")
			}
			for _, param := range v.Parameters {
				fmt.Printf("    - %s: %s\n", param.Name, param.Type)
			}
			fmt.Printf("  Line: %d, Column: %d\n", v.Line, v.Column)
			fmt.Printf("Size int: %d, float: %d", v.IntVarsCounter, v.FloatVarsCounter)
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
				case shared.Variable:
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

func (st *SymbolTable) GetGlobalScope() map[string]interface{} {
	return st.variables["global"]
}
