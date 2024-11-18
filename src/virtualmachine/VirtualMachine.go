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
	currFunction       string
}

func NewVirtualMachine(quads []shared.Quadruple, memManager *MemoryManager) *VirtualMachine {
	vm := &VirtualMachine{
		quads:              quads,
		memoryManager:      memManager,
		instructionPointer: 0,
		returnPointer:      shared.NewStack(),
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

	switch left := leftVal.(type) {
	case int:
		right, ok := rightVal.(int)
		if !ok {
			return fmt.Errorf("type mismatch: cannot perform integer operation with %T", rightVal)
		}

		switch quad.Operator {
		case "+":
			result = left + right
		case "-":
			result = left - right
		case "*":
			result = left * right
		case "/":
			if right == 0 {
				return fmt.Errorf("division by zero")
			}
			result = left / right
		}

	case float64:
		var right float64
		switch r := rightVal.(type) {
		case float64:
			right = r
		case int:
			right = float64(r)
		default:
			return fmt.Errorf("type mismatch: cannot perform float operation with %T", rightVal)
		}

		switch quad.Operator {
		case "+":
			result = left + right
		case "-":
			result = left - right
		case "*":
			result = left * right
		case "/":
			if right == 0 {
				return fmt.Errorf("division by zero")
			}
			result = left / right
		}
	}
	// Store result in memory
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
			fmt.Print(v, " ")
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
	value, err := vm.memoryManager.Load(valueAddr)
	if err != nil {
		return err
	}

	function, exists := vm.Functions[vm.currFunction]
	if exists {
		index := quad.RightOp.(int)
		currParam := function.Parameters[index]
		currParamAddr := currParam.Address
		if err := vm.memoryManager.Store(currParamAddr, value); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("function does not exist")
	}
	return nil
}

func (vm *VirtualMachine) executeEra(quad shared.Quadruple) error {
	functionName := quad.LeftOp.(string)
	vm.currFunction = functionName
	functionInfo, exists := vm.Functions[functionName]
	if exists {
		vm.memoryManager.PushNewFunctionSegment(false, functionInfo.IntVarsCount, functionInfo.FloatVarsCount)
	} else {
		return fmt.Errorf("function does not exist")
	}
	return nil
}

func (vm *VirtualMachine) executeEndproc(quad shared.Quadruple) error {
	vm.instructionPointer = vm.returnPointer.Pop().(int)
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
