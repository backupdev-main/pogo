package semantic

import (
	"fmt"
	"pogo/src/shared"
	"pogo/src/virtualmachine"
)

// Quadruple struct

type QuadrupleList struct {
	Quads         []shared.Quadruple
	OperatorStack *shared.Stack
	OperandStack  *shared.Stack
	TypeStack     *shared.Stack
	JumpStack     *shared.Stack
	TempCounter   int
	SemanticCube  *SemanticCube
	MemoryManager *virtualmachine.MemoryManager
}

func NewQuadrupleList() *QuadrupleList {
	return &QuadrupleList{
		Quads:         make([]shared.Quadruple, 0),
		OperatorStack: shared.NewStack(),
		OperandStack:  shared.NewStack(),
		TypeStack:     shared.NewStack(),
		JumpStack:     shared.NewStack(),
		TempCounter:   0,
		SemanticCube:  NewSemanticCube(),
		MemoryManager: virtualmachine.NewMemoryManager(),
	}

}

func (ql *QuadrupleList) HandleProgramStart() {
	quad := shared.Quadruple{
		Operator: "goto",
		LeftOp:   nil,
		RightOp:  nil,
		Result:   nil,
	}

	ql.Quads = append(ql.Quads, quad)
}

func (ql *QuadrupleList) NewTemp(tempType shared.Type) (int, error) {
	addr, err := ql.MemoryManager.AllocateTemp(tempType)
	if err != nil {
		return -1, err
	}
	return addr, nil
}

func (ql *QuadrupleList) HandleOpenParen() {
	ql.OperatorStack.Push("(")
}

func (ql *QuadrupleList) HandleCloseParen() error {
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
		rightType := ql.TypeStack.Pop().(shared.Type)
		leftOp := ql.OperandStack.Pop()
		leftType := ql.TypeStack.Pop().(shared.Type)

		resultType := ql.SemanticCube.GetResultType(leftType, rightType, operator)
		if resultType == shared.TypeError {
			return fmt.Errorf("type mismatch for operation %v %s %v", leftType, operator, rightType)
		}

		result, err := ql.NewTemp(resultType)

		if err != nil {
			return err
		}

		quad := shared.Quadruple{
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
		rightType := ql.TypeStack.Pop().(shared.Type)
		left := ql.OperandStack.Pop()
		leftType := ql.TypeStack.Pop().(shared.Type)
		op := ql.OperatorStack.Pop().(string)

		resultType := ql.SemanticCube.GetResultType(leftType, rightType, op)
		if resultType == shared.TypeError {
			return fmt.Errorf("type mismatch for operation %v %s %v", leftType, op, rightType)
		}

		result, err := ql.MemoryManager.AllocateTemp(resultType)

		if err != nil {
			return err
		}
		ql.TempCounter++

		ql.Quads = append(ql.Quads, shared.Quadruple{
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

func (ql *QuadrupleList) HandleFactor(value string, valueType shared.Type, symbolTable *SymbolTable) error {
	var addr int
	var err error

	if isNumeric(value) {
		addr, err = ql.MemoryManager.AllocateConstant(value)
		// fmt.Println("Numeric", addr, err)
	} else {
		addr, err = symbolTable.GetVariableAddress(value)
	}

	if err != nil {
		return fmt.Errorf("error allocating value: %v", err, value)
	}

	ql.OperandStack.Push(addr)
	ql.TypeStack.Push(valueType)
	//fmt.Printf("After HandleFactor(%v): Operands=%v, Operators=%v\n",
	//	value,
	//	ql.OperandStack,
	//	ql.OperatorStack)
	return nil
}

func (ql *QuadrupleList) HandleAssignment(target int, targetType shared.Type) error {
	if ql.OperandStack.Top() == nil {
		return fmt.Errorf("missing expression for assignment")
	}

	value := ql.OperandStack.Pop()
	valueType := ql.TypeStack.Pop().(shared.Type)

	// Check if assignment is valid using semantic cube
	resultType := ql.SemanticCube.GetResultType(targetType, valueType, "=")
	if resultType == shared.TypeError {
		return fmt.Errorf("cannot assign value of type %v to variable of type %v", valueType, targetType)
	}

	ql.Quads = append(ql.Quads, shared.Quadruple{
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

	if condType != shared.TypeInt {
		return fmt.Errorf("condition must be boolean (result of comparison)")
	}

	quad := shared.Quadruple{
		Operator: "gotof",
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

	ql.Quads = append(ql.Quads, shared.Quadruple{
		Operator: "goto",
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

	if condType != shared.TypeInt {
		return fmt.Errorf("condition must be 0 or 1 (result of comparison)")
	}

	quad := shared.Quadruple{
		Operator: "gotof",
		LeftOp:   condition,
		RightOp:  nil,
		Result:   nil,
	}

	jumpIndex := len(ql.Quads)
	ql.Quads = append(ql.Quads, quad)

	ql.JumpStack.Push(jumpIndex)

	return nil
}

func (ql *QuadrupleList) HandleElse() error {

	quad := shared.Quadruple{
		Operator: "goto",
		LeftOp:   nil,
		RightOp:  nil,
		Result:   nil,
	}

	gotoIndex := len(ql.Quads)
	ql.Quads = append(ql.Quads, quad)

	if ql.JumpStack.IsEmpty() {
		return fmt.Errorf("mismatched if-else: no corresponding if statement found")
	}
	falseJumpIndex := ql.JumpStack.Pop().(int)

	ql.Quads[falseJumpIndex].Result = len(ql.Quads)

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
	if v, ok := value.(string); ok {
		addr, err := ql.MemoryManager.GetStringAddress(v)
		if err != nil {
			return err
		}
		value = addr // Assign back to original value variable if needed
	}

	quad := shared.Quadruple{
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

func (ql *QuadrupleList) HandleERA(functionName string) error {
	quad := shared.Quadruple{
		Operator: "ERA",
		LeftOp:   functionName,
		RightOp:  nil,
		Result:   nil,
	}
	ql.Quads = append(ql.Quads, quad)
	return nil
}

func (ql *QuadrupleList) HandleParam(value interface{}, paramNum int) error {
	quad := shared.Quadruple{
		Operator: "PARAM",
		LeftOp:   value,
		RightOp:  paramNum,
		Result:   nil,
	}
	ql.Quads = append(ql.Quads, quad)
	return nil
}

func (ql *QuadrupleList) HandleGOSUB(functionName string, startQuad int) error {
	quad := shared.Quadruple{
		Operator: "GOSUB",
		LeftOp:   functionName,
		RightOp:  len(ql.Quads) + 1,
		Result:   startQuad,
	}
	ql.Quads = append(ql.Quads, quad)
	return nil
}

func (ql *QuadrupleList) HandleENDPROC() error {
	quad := shared.Quadruple{
		Operator: "ENDPROC",
		LeftOp:   nil,
		RightOp:  nil,
		Result:   nil,
	}
	ql.Quads = append(ql.Quads, quad)
	return nil
}

func (ql *QuadrupleList) PrintStacks() {
	fmt.Println("Stack Operators", ql.OperatorStack)
	fmt.Println("Stack Operands", ql.OperandStack)
	fmt.Println("Stack Types", ql.TypeStack)
}
