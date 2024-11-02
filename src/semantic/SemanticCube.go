package semantic

type SemanticCube struct {
	cube map[Type]map[Type]map[string]Type
}

// NewSemanticCube initializes the semantic cube for type checking
func NewSemanticCube() *SemanticCube {
	cube := &SemanticCube{
		cube: make(map[Type]map[Type]map[string]Type),
	}

	// Initialize semantic cube
	for _, t1 := range []Type{TypeInt, TypeFloat} {
		cube.cube[t1] = make(map[Type]map[string]Type)
		for _, t2 := range []Type{TypeInt, TypeFloat} {
			cube.cube[t1][t2] = make(map[string]Type)
		}
	}

	// Handle arithmetic operators
	arithOps := []string{"+", "-", "*", "/"}
	for _, op := range arithOps {
		// Int operations
		cube.cube[TypeInt][TypeInt][op] = TypeInt
		if op == "/" {
			cube.cube[TypeInt][TypeInt][op] = TypeFloat
		}
		cube.cube[TypeInt][TypeFloat][op] = TypeFloat

		cube.cube[TypeFloat][TypeInt][op] = TypeFloat
		cube.cube[TypeFloat][TypeFloat][op] = TypeFloat
	}

	relOps := []string{"<", ">", "==", "!=", "<=", ">="}
	for _, op := range relOps {
		cube.cube[TypeInt][TypeInt][op] = TypeInt
		cube.cube[TypeInt][TypeFloat][op] = TypeInt
		cube.cube[TypeFloat][TypeInt][op] = TypeInt
		cube.cube[TypeFloat][TypeFloat][op] = TypeInt
	}

	cube.cube[TypeFloat][TypeInt]["="] = TypeFloat
	cube.cube[TypeFloat][TypeFloat]["="] = TypeFloat
	cube.cube[TypeInt][TypeInt]["="] = TypeInt
	cube.cube[TypeInt][TypeFloat]["="] = TypeError

	return cube
}

func (sc *SemanticCube) GetResultType(t1, t2 Type, operator string) Type {
	// String operations are not allowed
	if t1 == TypeString || t2 == TypeString {
		return TypeError
	}

	if t1 == TypeError || t2 == TypeError {
		return TypeError
	}

	if result, exists := sc.cube[t1][t2][operator]; exists {
		return result
	}

	return TypeError
}

func (sc *SemanticCube) ValidatePrintItem(t Type) bool {
	return t == TypeInt || t == TypeFloat || t == TypeString
}
