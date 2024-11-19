package semantic

import "pogo/src/shared"

type SemanticCube struct {
	cube map[shared.Type]map[shared.Type]map[string]shared.Type
}

func NewSemanticCube() *SemanticCube {

	cube := &SemanticCube{
		cube: make(map[shared.Type]map[shared.Type]map[string]shared.Type),
	}
	// Initialize semantic cube
	for _, t1 := range []shared.Type{shared.TypeInt, shared.TypeFloat} {
		cube.cube[t1] = make(map[shared.Type]map[string]shared.Type)
		for _, t2 := range []shared.Type{shared.TypeInt, shared.TypeFloat} {
			cube.cube[t1][t2] = make(map[string]shared.Type)
		}
	}

	arithOps := []string{"+", "-", "*", "/"}
	for _, op := range arithOps {
		// Int operations
		cube.cube[shared.TypeInt][shared.TypeInt][op] = shared.TypeInt
		//if op == "/" {
		//	cube.cube[shared.TypeInt][shared.TypeInt][op] = shared.TypeFloat
		//}
		cube.cube[shared.TypeInt][shared.TypeFloat][op] = shared.TypeFloat

		cube.cube[shared.TypeFloat][shared.TypeInt][op] = shared.TypeFloat
		cube.cube[shared.TypeFloat][shared.TypeFloat][op] = shared.TypeFloat
	}

	relOps := []string{"<", ">", "==", "!=", "<=", ">="}
	for _, op := range relOps {
		cube.cube[shared.TypeInt][shared.TypeInt][op] = shared.TypeInt
		cube.cube[shared.TypeInt][shared.TypeFloat][op] = shared.TypeInt
		cube.cube[shared.TypeFloat][shared.TypeInt][op] = shared.TypeInt
		cube.cube[shared.TypeFloat][shared.TypeFloat][op] = shared.TypeInt
	}

	cube.cube[shared.TypeFloat][shared.TypeInt]["="] = shared.TypeFloat
	cube.cube[shared.TypeFloat][shared.TypeFloat]["="] = shared.TypeFloat
	cube.cube[shared.TypeInt][shared.TypeInt]["="] = shared.TypeInt
	cube.cube[shared.TypeInt][shared.TypeFloat]["="] = shared.TypeError

	return cube
}

func (sc *SemanticCube) GetResultType(t1, t2 shared.Type, operator string) shared.Type {
	if t1 == shared.TypeString || t2 == shared.TypeString {
		return shared.TypeError
	}

	if t1 == shared.TypeError || t2 == shared.TypeError {
		return shared.TypeError
	}

	if result, exists := sc.cube[t1][t2][operator]; exists {
		return result
	}

	return shared.TypeError
}

func (sc *SemanticCube) ValidatePrintItem(t shared.Type) bool {
	return t == shared.TypeInt || t == shared.TypeFloat || t == shared.TypeString
}
