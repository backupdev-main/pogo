package tests

//func TestSemanticCube(t *testing.T) {
//	cube := semantic.NewSemanticCube()
//
//	// Test cases structure: {type1, type2, operator, expectedResult}
//	testCases := []struct {
//		t1       shared.Type
//		t2       shared.Type
//		op       string
//		expected shared.Type
//	}{
//		// Arithmetic operations with integers
//		{shared.TypeInt, shared.TypeInt, "+", shared.TypeInt},
//		{shared.TypeInt, shared.TypeInt, "-", shared.TypeInt},
//		{shared.TypeInt, shared.TypeInt, "*", shared.TypeInt},
//		{shared.TypeInt, shared.TypeInt, "/", shared.TypeInt},
//
//		// Arithmetic operations with mixed shared (should always return float)
//		{shared.TypeInt, shared.TypeFloat, "+", shared.TypeFloat},
//		{shared.TypeFloat, shared.TypeInt, "+", shared.TypeFloat},
//		{shared.TypeInt, shared.TypeFloat, "-", shared.TypeFloat},
//		{shared.TypeFloat, shared.TypeInt, "-", shared.TypeFloat},
//		{shared.TypeInt, shared.TypeFloat, "*", shared.TypeFloat},
//		{shared.TypeFloat, shared.TypeInt, "*", shared.TypeFloat},
//		{shared.TypeInt, shared.TypeFloat, "/", shared.TypeFloat},
//		{shared.TypeFloat, shared.TypeInt, "/", shared.TypeFloat},
//
//		// Arithmetic operations with floats
//		{shared.TypeFloat, shared.TypeFloat, "+", shared.TypeFloat},
//		{shared.TypeFloat, shared.TypeFloat, "-", shared.TypeFloat},
//		{shared.TypeFloat, shared.TypeFloat, "*", shared.TypeFloat},
//		{shared.TypeFloat, shared.TypeFloat, "/", shared.TypeFloat},
//
//		// Relational operations with integers
//		{shared.TypeInt, shared.TypeInt, "<", shared.TypeInt},
//		{shared.TypeInt, shared.TypeInt, ">", shared.TypeInt},
//		{shared.TypeInt, shared.TypeInt, "<=", semantic.TypeInt},
//		{semantic.TypeInt, semantic.TypeInt, ">=", semantic.TypeInt},
//		{semantic.TypeInt, semantic.TypeInt, "==", semantic.TypeInt},
//		{semantic.TypeInt, semantic.TypeInt, "!=", semantic.TypeInt},
//
//		// Relational operations with mixed shared (should return int as boolean)
//		{semantic.TypeInt, semantic.TypeFloat, "<", semantic.TypeInt},
//		{semantic.TypeFloat, semantic.TypeInt, ">", semantic.TypeInt},
//		{semantic.TypeInt, semantic.TypeFloat, "<=", semantic.TypeInt},
//		{semantic.TypeFloat, semantic.TypeInt, ">=", semantic.TypeInt},
//		{semantic.TypeInt, semantic.TypeFloat, "==", semantic.TypeInt},
//		{semantic.TypeFloat, semantic.TypeInt, "!=", semantic.TypeInt},
//
//		// Relational operations with floats
//		{semantic.TypeFloat, semantic.TypeFloat, "<", semantic.TypeInt},
//		{semantic.TypeFloat, semantic.TypeFloat, ">", semantic.TypeInt},
//		{semantic.TypeFloat, semantic.TypeFloat, "<=", semantic.TypeInt},
//		{semantic.TypeFloat, semantic.TypeFloat, ">=", semantic.TypeInt},
//		{semantic.TypeFloat, semantic.TypeFloat, "==", semantic.TypeInt},
//		{semantic.TypeFloat, semantic.TypeFloat, "!=", semantic.TypeInt},
//
//		// Error cases
//		{semantic.TypeError, semantic.TypeInt, "+", semantic.TypeError},
//		{semantic.TypeInt, semantic.TypeError, "+", semantic.TypeError},
//		{semantic.TypeError, semantic.TypeError, "+", semantic.TypeError},
//	}
//
//	for i, tc := range testCases {
//		result := cube.GetResultType(tc.t1, tc.t2, tc.op)
//		if result != tc.expected {
//			t.Errorf("Test case %d failed: %s %s %s = %s, expected %s",
//				i+1,
//				tc.t1.String(),
//				tc.op,
//				tc.t2.String(),
//				result.String(),
//				tc.expected.String())
//		}
//	}
//}
//
//func TestSemanticCubeOperators(t *testing.T) {
//	cube := semantic.NewSemanticCube()
//	shared := []semantic.Type{semantic.TypeInt, semantic.TypeFloat}
//	arithmeticOps := []string{"+", "-", "*", "/"}
//	relationalOps := []string{"<", ">", "<=", ">=", "==", "!="}
//
//	// Test arithmetic operators
//	for _, t1 := range shared {
//		for _, t2 := range shared {
//			for _, op := range arithmeticOps {
//				result := cube.GetResultType(t1, t2, op)
//				if t1 == semantic.TypeInt && t2 == semantic.TypeInt {
//					if result != semantic.TypeInt {
//						t.Errorf("Arithmetic: %s %s %s should be int, got %s",
//							t1.String(), op, t2.String(), result.String())
//					}
//				} else {
//					if result != semantic.TypeFloat {
//						t.Errorf("Arithmetic: %s %s %s should be float, got %s",
//							t1.String(), op, t2.String(), result.String())
//					}
//				}
//			}
//		}
//	}
//
//	// Test relational operators
//	for _, t1 := range shared {
//		for _, t2 := range shared {
//			for _, op := range relationalOps {
//				result := cube.GetResultType(t1, t2, op)
//				if result != semantic.TypeInt {
//					t.Errorf("Relational: %s %s %s should return int (boolean), got %s",
//						t1.String(), op, t2.String(), result.String())
//				}
//			}
//		}
//	}
//}
//
//func TestSymbolTableStringVariables(t *testing.T) {
//	st := semantic.NewSymbolTable()
//
//	// Test that string variables cannot be declared
//	err := st.AddVariable("str", semantic.TypeString, 1, 1)
//	if err == nil {
//		t.Error("Expected error when declaring string variable, got nil")
//	}
//}
//
//func TestTypeString(t *testing.T) {
//	testCases := []struct {
//		t        semantic.Type
//		expected string
//	}{
//		{semantic.TypeInt, "int"},
//		{semantic.TypeFloat, "float"},
//		{semantic.TypeError, "error"},
//		{semantic.Type(999), "error"},
//	}
//
//	for _, tc := range testCases {
//		if tc.t.String() != tc.expected {
//			t.Errorf("Type %d string representation should be %s, got %s",
//				tc.t, tc.expected, tc.t.String())
//		}
//	}
//}
