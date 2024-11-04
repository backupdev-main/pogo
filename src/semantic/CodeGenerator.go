package semantic

import "fmt"

// Quadruple struct
type Quadruple struct {
	Operator string      // The operation to be performed
	LeftOp   interface{} // Left operand
	RightOp  interface{} // Right operand
	Result   interface{} // Where the result will be stored
}

type QuadrupleList struct {
	Quads         []Quadruple
	OperatorStack *Stack
	OperandStack  *Stack
	TypeStack     *Stack
	JumpStack     *Stack
	TempCounter   int
	SemanticCube  *SemanticCube
}

func NewQuadrupleList() *QuadrupleList {
	return &QuadrupleList{
		Quads:         make([]Quadruple, 0),
		OperatorStack: NewStack(),
		OperandStack:  NewStack(),
		TypeStack:     NewStack(),
		JumpStack:     NewStack(),
		TempCounter:   0,
		SemanticCube:  NewSemanticCube(),
	}
}

func (ql *QuadrupleList) NewTemp() string {
	temp := fmt.Sprintf("t%d", ql.TempCounter)
	ql.TempCounter++
	return temp
}

func (ql *QuadrupleList) HandleOpenParen() {
	ql.OperatorStack.Push("(")
}

func (ql *QuadrupleList) HandleCloseParen() error {
	// Pop and process all operators until we find the matching open parenthesis
	for !ql.OperatorStack.IsEmpty() {
		topOp := ql.OperatorStack.Top()
		if topOp == nil {
			return fmt.Errorf("mismatched parentheses: no opening parenthesis found")
		}

		if topOp.(string) == "(" {
			ql.OperatorStack.Pop()
			return nil
		}

		operator := ql.OperatorStack.Pop().(string)

		if ql.OperandStack.Size() < 2 {
			return fmt.Errorf("insufficient operands for operator %s", operator)
		}

		rightOp := ql.OperandStack.Pop()
		rightType := ql.TypeStack.Pop().(Type)
		leftOp := ql.OperandStack.Pop()
		leftType := ql.TypeStack.Pop().(Type)

		resultType := ql.SemanticCube.GetResultType(leftType, rightType, operator)
		if resultType == TypeError {
			return fmt.Errorf("type mismatch for operation %v %s %v", leftType, operator, rightType)
		}

		result := ql.NewTemp()

		quad := Quadruple{
			Operator: operator,
			LeftOp:   leftOp,
			RightOp:  rightOp,
			Result:   result,
		}
		ql.Quads = append(ql.Quads, quad)

		ql.OperandStack.Push(result)
		ql.TypeStack.Push(resultType)
	}

	return fmt.Errorf("mismatched parentheses: no opening parenthesis found")
}

func (ql *QuadrupleList) HandleOp() error {
	if ql.OperatorStack.Top() != nil {
		right := ql.OperandStack.Pop()
		rightType := ql.TypeStack.Pop().(Type)
		left := ql.OperandStack.Pop()
		leftType := ql.TypeStack.Pop().(Type)
		op := ql.OperatorStack.Pop().(string)

		resultType := ql.SemanticCube.GetResultType(leftType, rightType, op)
		if resultType == TypeError {
			return fmt.Errorf("type mismatch for operation %v %s %v", leftType, op, rightType)
		}

		result := fmt.Sprintf("t%d", ql.TempCounter)
		ql.TempCounter++

		ql.Quads = append(ql.Quads, Quadruple{
			Operator: op,
			LeftOp:   left,
			RightOp:  right,
			Result:   result,
		})

		ql.OperandStack.Push(result)
		ql.TypeStack.Push(resultType)
	}

	return nil
}

func (ql *QuadrupleList) HandleFactor(value interface{}, valueType Type) {
	ql.OperandStack.Push(value)
	ql.TypeStack.Push(valueType)
	//fmt.Printf("After HandleFactor(%v): Operands=%v, Operators=%v\n",
	//	value,
	//	ql.OperandStack,
	//	ql.OperatorStack)
}

func (ql *QuadrupleList) HandleAssignment(target string, targetType Type) error {
	if ql.OperandStack.Top() == nil {
		return fmt.Errorf("missing expression for assignment")
	}

	value := ql.OperandStack.Pop()
	valueType := ql.TypeStack.Pop().(Type)

	// Check if assignment is valid using semantic cube
	resultType := ql.SemanticCube.GetResultType(targetType, valueType, "=")
	if resultType == TypeError {
		return fmt.Errorf("cannot assign value of type %v to variable of type %v", valueType, targetType)
	}

	ql.Quads = append(ql.Quads, Quadruple{
		Operator: "=",
		LeftOp:   value,
		RightOp:  nil,
		Result:   target,
	})

	return nil
}

func (ql *QuadrupleList) HandleWhileStart() int {
	return len(ql.Quads)
}

func (ql *QuadrupleList) HandleWhileCondition() error {
	if ql.OperandStack.IsEmpty() {
		return fmt.Errorf("missing condition for while statement")
	}

	condition := ql.OperandStack.Pop()
	condType := ql.TypeStack.Pop()

	if condType != TypeInt {
		return fmt.Errorf("condition must be boolean (result of comparison)")
	}

	quad := Quadruple{
		Operator: "GotoF",
		LeftOp:   condition,
		RightOp:  nil,
		Result:   nil,
	}

	jumpIndex := len(ql.Quads)
	ql.Quads = append(ql.Quads, quad)

	ql.JumpStack.Push(jumpIndex)

	return nil
}

func (ql *QuadrupleList) HandleWhileEnd(startIndex int) error {
	if ql.JumpStack.IsEmpty() {
		return fmt.Errorf("mismatched while: no pending jumps found")
	}

	ql.Quads = append(ql.Quads, Quadruple{
		Operator: "Goto",
		LeftOp:   nil,
		RightOp:  nil,
		Result:   startIndex,
	})

	falseJumpIndex := ql.JumpStack.Pop().(int)
	ql.Quads[falseJumpIndex].Result = len(ql.Quads)

	return nil
}

func (ql *QuadrupleList) HandleIfStatement() error {
	if ql.OperandStack.IsEmpty() {
		return fmt.Errorf("missing condition for if statement")
	}

	condition := ql.OperandStack.Pop()
	condType := ql.TypeStack.Pop()

	if condType != TypeInt {
		return fmt.Errorf("condition must be boolean (result of comparison)")
	}

	quad := Quadruple{
		Operator: "GotoF",
		LeftOp:   condition,
		RightOp:  nil,
		Result:   nil, // This will be filled in later
	}

	jumpIndex := len(ql.Quads)
	ql.Quads = append(ql.Quads, quad)

	ql.JumpStack.Push(jumpIndex)

	return nil
}

func (ql *QuadrupleList) HandleElse() error {

	quad := Quadruple{
		Operator: "Goto",
		LeftOp:   nil,
		RightOp:  nil,
		Result:   nil,
	}

	// Add the Goto quadruple
	gotoIndex := len(ql.Quads)
	ql.Quads = append(ql.Quads, quad)

	// Fill the previous GotoF (from HandleIf)
	if ql.JumpStack.IsEmpty() {
		return fmt.Errorf("mismatched if-else: no corresponding if statement found")
	}
	falseJumpIndex := ql.JumpStack.Pop().(int)

	// The false jump should point to the next quadruple
	ql.Quads[falseJumpIndex].Result = len(ql.Quads)

	// Push the Goto index for later backpatching
	ql.JumpStack.Push(gotoIndex)

	return nil
}

func (ql *QuadrupleList) HandleEndIf() error {
	if ql.JumpStack.IsEmpty() {
		return fmt.Errorf("mismatched if-else: no pending jumps found")
	}

	jumpIndex := ql.JumpStack.Pop().(int)

	ql.Quads[jumpIndex].Result = len(ql.Quads)

	return nil
}

func (ql *QuadrupleList) HandlePrint(value interface{}) error {
	quad := Quadruple{
		Operator: "print",
		LeftOp:   value,
		RightOp:  nil,
		Result:   nil,
	}

	ql.Quads = append(ql.Quads, quad)
	return nil
}

func (ql *QuadrupleList) Print() {
	fmt.Println("Generated Quadruples:")
	for i, quad := range ql.Quads {
		fmt.Printf("%d: (%v, %v, %v, %v)\n", i, quad.Operator, quad.LeftOp, quad.RightOp, quad.Result)
	}
}

func (ql *QuadrupleList) PrintStacks() {
	fmt.Println("Stack Operators", ql.OperatorStack)
	fmt.Println("Stack Operands", ql.OperandStack)
	fmt.Println("Stack Types", ql.TypeStack)
}
