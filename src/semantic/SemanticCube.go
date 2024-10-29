package semantic

type SemanticCube struct {
	cube map[Type]map[Type]map[string]Type
}

// NewSemanticCube initializes the semantic cube for type checking
func NewSemanticCube() *SemanticCube {
	cube := &SemanticCube{
		cube: make(map[Type]map[Type]map[string]Type),
	}

	// Initialize the cube for numeric types only
	for _, t1 := range []Type{TypeInt, TypeFloat} {
		cube.cube[t1] = make(map[Type]map[string]Type)
		for _, t2 := range []Type{TypeInt, TypeFloat} {
			cube.cube[t1][t2] = make(map[string]Type)

			// Arithmetic operators (+, -, *, /)
			for _, op := range []string{"+", "-", "*", "/"} {
				if t1 == TypeInt && t2 == TypeInt {
					cube.cube[t1][t2][op] = TypeInt
				} else {
					cube.cube[t1][t2][op] = TypeFloat
				}
			}

			// Relational operators (<, >, <=, >=, ==, !=)
			for _, op := range []string{"<", ">", "<=", ">=", "==", "!="} {
				cube.cube[t1][t2][op] = TypeInt // Using int as boolean (0/1)
			}

			if t1 == t2 {
				cube.cube[t1][t2]["="] = t1
			} else if t1 == TypeFloat && t2 == TypeInt {
				cube.cube[t1][t2]["="] = TypeFloat
			} else {
				cube.cube[t1][t2]["="] = TypeError
			}
		}
	}

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
