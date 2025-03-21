package virtualmachine

import (
	"fmt"
	"pogo/src/shared"
	"strings"
)

type VirtualMachine struct {
	quads              []shared.Quadruple
	memoryManager      *MemoryManager
	Functions          map[string]shared.FunctionInfo
	instructionPointer int
	returnPointer      *shared.Stack
	functionStack      *shared.Stack
}

func NewVirtualMachine(quads []shared.Quadruple, memManager *MemoryManager) *VirtualMachine {
	vm := &VirtualMachine{
		quads:              quads,
		memoryManager:      memManager,
		instructionPointer: 0,
		returnPointer:      shared.NewStack(),
		functionStack:      shared.NewStack(),
		Functions:          make(map[string]shared.FunctionInfo),
	}

	return vm
}

func (vm *VirtualMachine) Execute() error {
	//fmt.Println("These are the quads", vm.quads)
	// fmt.Println("These are the functions", vm.Functions)
	// vm.memoryManager.InitializeMemory()
	vm.memoryManager.InitializeMemory()
	for vm.instructionPointer < len(vm.quads) {
		quad := vm.quads[vm.instructionPointer]

		if err := vm.executeQuadruple(quad); err != nil {
			return fmt.Errorf("error at instruction %d: %v", vm.instructionPointer, err)
		}

		vm.instructionPointer++
	}
	return nil
}

func (vm *VirtualMachine) executeQuadruple(quad shared.Quadruple) error {
	switch quad.Operator {
	case "+", "-", "*", "/":
		// fmt.Println("Entering with", quad.Operator)
		return vm.executeArithmetic(quad)
	case "=":
		return vm.executeAssignment(quad)
	case "<", ">", "==", "!=":
		return vm.executeComparison(quad)
	case "print":
		return vm.executePrint(quad)
	case "goto":
		return vm.executeGoto(quad)
	case "gotof":
		return vm.executeGotoF(quad)
	case "gosub":
		return vm.executeGosub(quad)
	case "era":
		return vm.executeEra(quad)
	case "endproc":
		return vm.executeEndproc(quad)
	case "param":
		return vm.executeParam(quad)
	}

	return nil
}

func (vm *VirtualMachine) executeArithmetic(quad shared.Quadruple) error {
	leftVal, err := vm.memoryManager.Load(quad.LeftOp.(int))

	if err != nil {
		return fmt.Errorf("failed to load left operand: %v", err)
	}

	rightVal, err := vm.memoryManager.Load(quad.RightOp.(int))
	if err != nil {
		return fmt.Errorf("failed to load right operand: %v", err)
	}

	var result interface{}
	var leftFloat, rightFloat float64
	isFloatOperation := false

	switch l := leftVal.(type) {
	case int:
		leftFloat = float64(l)
	case float64:
		leftFloat = l
		isFloatOperation = true
	default:
		return fmt.Errorf("invalid left operand type: %T", leftVal)
	}

	switch r := rightVal.(type) {
	case int:
		rightFloat = float64(r)
	case float64:
		rightFloat = r
		isFloatOperation = true
	default:
		return fmt.Errorf("invalid right operand type: %T", rightVal)
	}

	var floatResult float64
	switch quad.Operator {
	case "+":
		floatResult = leftFloat + rightFloat
	case "-":
		floatResult = leftFloat - rightFloat
	case "*":
		floatResult = leftFloat * rightFloat
	case "/":
		if rightFloat == 0 {
			return fmt.Errorf("division by zero")
		}
		floatResult = leftFloat / rightFloat
		isFloatOperation = true // Division always returns float
	default:
		return fmt.Errorf("unknown arithmetic operator: %s", quad.Operator)
	}

	if isFloatOperation {
		result = floatResult
	} else {
		result = int(floatResult)
	}

	return vm.memoryManager.Store(quad.Result.(int), result)
}

func (vm *VirtualMachine) executeAssignment(quad shared.Quadruple) error {
	value, err := vm.memoryManager.Load(quad.LeftOp.(int))
	if err != nil {
		return fmt.Errorf("failed to load source value: %v", err)
	}

	return vm.memoryManager.Store(quad.Result.(int), value)
}

func (vm *VirtualMachine) executeComparison(quad shared.Quadruple) error {
	// fmt.Println("Entering execution")
	leftVal, err := vm.memoryManager.Load(quad.LeftOp.(int))

	if err != nil {
		return fmt.Errorf("failed to load left operand: %v", err)
	}

	rightVal, err := vm.memoryManager.Load(quad.RightOp.(int))
	if err != nil {
		return fmt.Errorf("failed to load right operand: %v", err)
	}
	// fmt.Println("This is the rightVal", rightVal)

	var leftFloat, rightFloat float64

	switch v := leftVal.(type) {
	case int:
		leftFloat = float64(v)
	case float64:
		leftFloat = v
	default:
		return fmt.Errorf("invalid type for comparison: %T", leftVal)
	}

	switch v := rightVal.(type) {
	case int:
		rightFloat = float64(v)
	case float64:
		rightFloat = v
	default:
		return fmt.Errorf("invalid type for comparison: %T", rightVal)
	}

	var result bool
	var intResult int

	switch quad.Operator {
	case "<":
		result = leftFloat < rightFloat
	case ">":
		result = leftFloat > rightFloat
	case "==":
		result = leftFloat == rightFloat
	case "!=":
		result = leftFloat != rightFloat
	}

	if result {
		intResult = 1
	} else {
		intResult = 0
	}
	// fmt.Println("This is where we store", quad.Result)
	return vm.memoryManager.Store(quad.Result.(int), intResult)
}

func (vm *VirtualMachine) executePrint(quad shared.Quadruple) error {
	items := quad.LeftOp.([]int)
	for i, item := range items {
		value, err := vm.memoryManager.Load(item)
		if err != nil {
			return fmt.Errorf("failed to load print value: %v", err)
		}

		switch v := value.(type) {
		case string:
			cleanStr := strings.Trim(v, "\"")
			fmt.Print(cleanStr, " ")
		case int:
			fmt.Print(v, " ")
		case float64:
			fmt.Printf("%.2f ", v)
		default:
			return fmt.Errorf("unsupported type for printing: %T", value)
		}

		if i == len(items)-1 {
			fmt.Println()
		}
	}
	return nil
}

func (vm *VirtualMachine) executeGoto(quad shared.Quadruple) error {
	vm.instructionPointer = quad.Result.(int) - 1
	return nil
}

func (vm *VirtualMachine) executeGotoF(quad shared.Quadruple) error {
	condValue, err := vm.memoryManager.Load(quad.LeftOp.(int))
	if err != nil {
		return err
	}

	if condValue == 0 {
		vm.instructionPointer = quad.Result.(int) - 1
	}

	return nil
}

func (vm *VirtualMachine) executeParam(quad shared.Quadruple) error {
	valueAddr := quad.LeftOp.(int)

	if err := vm.memoryManager.PopFunctionSegment(); err != nil {
		return err
	}
	currentFunction := vm.functionStack.Pop()

	value, err := vm.memoryManager.Load(valueAddr)
	if err != nil {
		return err
	}

	vm.functionStack.Push(currentFunction)
	function, exists := vm.Functions[currentFunction.(string)]
	if exists {
		vm.memoryManager.PushNewFunctionSegment(false, function.IntVarsCount, function.FloatVarsCount)
	} else {
		return fmt.Errorf("function does not exist")
	}

	index := quad.RightOp.(int)
	currParam := function.Parameters[index]
	currParamAddr := currParam.Address
	if err := vm.memoryManager.Store(currParamAddr, value); err != nil {
		return err
	}

	return nil
}

func (vm *VirtualMachine) executeEra(quad shared.Quadruple) error {
	functionName := quad.LeftOp.(string)
	vm.functionStack.Push(functionName)
	functionInfo, exists := vm.Functions[functionName]
	//fmt.Println("This is the functionInfo", functionInfo)
	//fmt.Println("These are the counts", functionInfo.IntVarsCount, functionInfo.FloatVarsCount)
	if exists {
		vm.memoryManager.PushNewFunctionSegment(false, functionInfo.IntVarsCount, functionInfo.FloatVarsCount)
	} else {
		return fmt.Errorf("function does not exist")
	}
	return nil
}

func (vm *VirtualMachine) executeEndproc(quad shared.Quadruple) error {
	vm.instructionPointer = vm.returnPointer.Pop().(int)
	vm.functionStack.Pop()
	if err := vm.memoryManager.PopFunctionSegment(); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachine) executeGosub(quad shared.Quadruple) error {
	start := quad.Result
	vm.returnPointer.Push(vm.instructionPointer)
	vm.instructionPointer = start.(int) - 1
	return nil
}
