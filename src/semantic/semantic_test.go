package semantic

import "testing"

func TestSemanticCube(t *testing.T) {
	cube := NewSemanticCube()

	// Test cases structure: {type1, type2, operator, expectedResult}
	testCases := []struct {
		t1       Type
		t2       Type
		op       string
		expected Type
	}{
		// Arithmetic operations with integers
		{TypeInt, TypeInt, "+", TypeInt},
		{TypeInt, TypeInt, "-", TypeInt},
		{TypeInt, TypeInt, "*", TypeInt},
		{TypeInt, TypeInt, "/", TypeInt},

		// Arithmetic operations with mixed types (should always return float)
		{TypeInt, TypeFloat, "+", TypeFloat},
		{TypeFloat, TypeInt, "+", TypeFloat},
		{TypeInt, TypeFloat, "-", TypeFloat},
		{TypeFloat, TypeInt, "-", TypeFloat},
		{TypeInt, TypeFloat, "*", TypeFloat},
		{TypeFloat, TypeInt, "*", TypeFloat},
		{TypeInt, TypeFloat, "/", TypeFloat},
		{TypeFloat, TypeInt, "/", TypeFloat},

		// Arithmetic operations with floats
		{TypeFloat, TypeFloat, "+", TypeFloat},
		{TypeFloat, TypeFloat, "-", TypeFloat},
		{TypeFloat, TypeFloat, "*", TypeFloat},
		{TypeFloat, TypeFloat, "/", TypeFloat},

		// Relational operations with integers
		{TypeInt, TypeInt, "<", TypeInt},
		{TypeInt, TypeInt, ">", TypeInt},
		{TypeInt, TypeInt, "<=", TypeInt},
		{TypeInt, TypeInt, ">=", TypeInt},
		{TypeInt, TypeInt, "==", TypeInt},
		{TypeInt, TypeInt, "!=", TypeInt},

		// Relational operations with mixed types (should return int as boolean)
		{TypeInt, TypeFloat, "<", TypeInt},
		{TypeFloat, TypeInt, ">", TypeInt},
		{TypeInt, TypeFloat, "<=", TypeInt},
		{TypeFloat, TypeInt, ">=", TypeInt},
		{TypeInt, TypeFloat, "==", TypeInt},
		{TypeFloat, TypeInt, "!=", TypeInt},

		// Relational operations with floats
		{TypeFloat, TypeFloat, "<", TypeInt},
		{TypeFloat, TypeFloat, ">", TypeInt},
		{TypeFloat, TypeFloat, "<=", TypeInt},
		{TypeFloat, TypeFloat, ">=", TypeInt},
		{TypeFloat, TypeFloat, "==", TypeInt},
		{TypeFloat, TypeFloat, "!=", TypeInt},

		// Error cases
		{TypeError, TypeInt, "+", TypeError},
		{TypeInt, TypeError, "+", TypeError},
		{TypeError, TypeError, "+", TypeError},
	}

	for i, tc := range testCases {
		result := cube.GetResultType(tc.t1, tc.t2, tc.op)
		if result != tc.expected {
			t.Errorf("Test case %d failed: %s %s %s = %s, expected %s",
				i+1,
				tc.t1.String(),
				tc.op,
				tc.t2.String(),
				result.String(),
				tc.expected.String())
		}
	}
}

func TestSemanticCubeOperators(t *testing.T) {
	cube := NewSemanticCube()
	types := []Type{TypeInt, TypeFloat}
	arithmeticOps := []string{"+", "-", "*", "/"}
	relationalOps := []string{"<", ">", "<=", ">=", "==", "!="}

	// Test arithmetic operators
	for _, t1 := range types {
		for _, t2 := range types {
			for _, op := range arithmeticOps {
				result := cube.GetResultType(t1, t2, op)
				if t1 == TypeInt && t2 == TypeInt {
					if result != TypeInt {
						t.Errorf("Arithmetic: %s %s %s should be int, got %s",
							t1.String(), op, t2.String(), result.String())
					}
				} else {
					if result != TypeFloat {
						t.Errorf("Arithmetic: %s %s %s should be float, got %s",
							t1.String(), op, t2.String(), result.String())
					}
				}
			}
		}
	}

	// Test relational operators
	for _, t1 := range types {
		for _, t2 := range types {
			for _, op := range relationalOps {
				result := cube.GetResultType(t1, t2, op)
				if result != TypeInt {
					t.Errorf("Relational: %s %s %s should return int (boolean), got %s",
						t1.String(), op, t2.String(), result.String())
				}
			}
		}
	}
}

func TestSymbolTableStringVariables(t *testing.T) {
	st := NewSymbolTable()

	// Test that string variables cannot be declared
	err := st.AddVariable("str", TypeString, 1, 1)
	if err == nil {
		t.Error("Expected error when declaring string variable, got nil")
	}
}

func TestTypeString(t *testing.T) {
	testCases := []struct {
		t        Type
		expected string
	}{
		{TypeInt, "int"},
		{TypeFloat, "float"},
		{TypeError, "error"},
		{Type(999), "error"},
	}

	for _, tc := range testCases {
		if tc.t.String() != tc.expected {
			t.Errorf("Type %d string representation should be %s, got %s",
				tc.t, tc.expected, tc.t.String())
		}
	}
}
