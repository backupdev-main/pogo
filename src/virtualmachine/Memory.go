package virtualmachine

import (
	"fmt"
	"pogo/src/shared"
	"strconv"
)

const (
	GLOBAL_START   = 0
	LOCAL_START    = 4000
	TEMP_START     = 8000
	CONSTANT_START = 12000

	GLOBAL_INT_START   = 0
	GLOBAL_INT_END     = 1999
	GLOBAL_FLOAT_START = 2000
	GLOBAL_FLOAT_END   = 3999

	LOCAL_INT_START   = 4000
	LOCAL_INT_END     = 5999
	LOCAL_FLOAT_START = 6000
	LOCAL_FLOAT_END   = 7999

	TEMP_INT_START   = 8000
	TEMP_INT_END     = 9999
	TEMP_FLOAT_START = 10000
	TEMP_FLOAT_END   = 11999

	CONSTANT_INT_START   = 12000
	CONSTANT_INT_END     = 12999
	CONSTANT_FLOAT_START = 13000
	CONSTANT_FLOAT_END   = 13999
	CONSTANT_STR_START   = 14000
	CONSTANT_STR_END     = 15999

	MEMORY_SEGMENT_SIZE = 4000
	TOTAL_MEMORY_SIZE   = 16000
)

type FunctionMemorySegment struct {
	localMemory    []interface{}
	localIntPtr    int
	localFloatPtr  int
	intVarsCount   int
	floatVarsCount int
}

type MemoryManager struct {
	globalMemory   [MEMORY_SEGMENT_SIZE]interface{}
	globalIntPtr   int
	globalFloatPtr int

	tempMemory   [MEMORY_SEGMENT_SIZE]interface{}
	tempIntPtr   int
	tempFloatPtr int

	ConstantMemory   [MEMORY_SEGMENT_SIZE]interface{}
	constantIntPtr   int
	constantFloatPtr int
	constantStrPtr   int

	// map for constant reusing / not restoring the same constant
	constantMap       map[string]int
	constantStringMap map[string]int

	memoryStack    []FunctionMemorySegment
	currentSegment *FunctionMemorySegment
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		globalIntPtr:      GLOBAL_INT_START,
		globalFloatPtr:    GLOBAL_FLOAT_START,
		tempIntPtr:        TEMP_INT_START,
		tempFloatPtr:      TEMP_FLOAT_START,
		constantIntPtr:    CONSTANT_INT_START,
		constantFloatPtr:  CONSTANT_FLOAT_START,
		constantStrPtr:    CONSTANT_STR_START,
		constantMap:       make(map[string]int),
		constantStringMap: make(map[string]int),
		memoryStack:       make([]FunctionMemorySegment, 0),
		currentSegment:    nil,
	}
}

func (mm *MemoryManager) AllocateGlobal(varType shared.Type) (int, error) {
	switch varType {
	case shared.TypeInt:
		if mm.globalIntPtr >= GLOBAL_INT_END {
			return -1, fmt.Errorf("global integer memory overflow")
		}
		addr := mm.globalIntPtr
		mm.globalIntPtr++
		return addr, nil
	case shared.TypeFloat:
		if mm.globalFloatPtr >= GLOBAL_FLOAT_END {
			return -1, fmt.Errorf("global float memory overflow")
		}
		addr := mm.globalFloatPtr
		mm.globalFloatPtr++
		return addr, nil
	default:
		return -1, fmt.Errorf("unsupported type for global allocation")
	}
}

func (mm *MemoryManager) AllocateTemp(varType shared.Type) (int, error) {
	switch varType {
	case shared.TypeInt:
		if mm.tempIntPtr >= TEMP_INT_END {
			return -1, fmt.Errorf("temporary integer memory overflow")
		}
		addr := mm.tempIntPtr
		mm.tempIntPtr++
		return addr, nil
	case shared.TypeFloat:
		if mm.tempFloatPtr >= TEMP_FLOAT_END {
			return -1, fmt.Errorf("temporary float memory overflow")
		}
		addr := mm.tempFloatPtr
		mm.tempFloatPtr++
		return addr, nil
	default:
		return -1, fmt.Errorf("unsupported type for temporary allocation")
	}
}

func (mm *MemoryManager) AllocateConstant(value string) (int, error) {

	if addr, exists := mm.constantMap[value]; exists {
		return addr, nil
	}

	if intVal, ok := strconv.Atoi(value); ok == nil {
		if mm.constantIntPtr >= CONSTANT_INT_END {
			return -1, fmt.Errorf("constant integer memory overflow")
		}
		addr := mm.constantIntPtr
		// fmt.Println("Allocating value ", intVal, "at ", addr-CONSTANT_START)
		mm.ConstantMemory[addr-CONSTANT_START] = intVal
		mm.constantMap[value] = addr
		mm.constantIntPtr++
		return addr, nil
	}

	if floatVal, ok := strconv.ParseFloat(value, 64); ok == nil {
		if mm.constantFloatPtr >= CONSTANT_FLOAT_END {
			return -1, fmt.Errorf("constant integer memory overflow")
		}
		addr := mm.constantFloatPtr
		mm.ConstantMemory[addr-CONSTANT_START] = floatVal
		mm.constantMap[value] = addr
		mm.constantFloatPtr++
		return addr, nil
	}
	fmt.Println("Returning -1 for", value)
	return -1, fmt.Errorf("invalid constant value: %s", value)
}

func (mm *MemoryManager) AllocateLocal(varType shared.Type) (int, error) {
	switch varType {
	case shared.TypeInt:
		if mm.currentSegment.localIntPtr >= LOCAL_INT_END {
			return -1, fmt.Errorf("local integer memory overflow")
		}
		addr := mm.currentSegment.localIntPtr
		mm.currentSegment.localIntPtr++
		return addr, nil

	case shared.TypeFloat:
		if mm.currentSegment.localFloatPtr >= LOCAL_FLOAT_END {
			return -1, fmt.Errorf("local float memory overflow")
		}
		addr := mm.currentSegment.localFloatPtr
		mm.currentSegment.localFloatPtr++
		return addr, nil

	default:
		return -1, fmt.Errorf("unsupported type for local allocation")
	}
}

func (mm *MemoryManager) GetStringAddress(value string) (int, error) {
	if addr, exists := mm.constantStringMap[value]; exists {
		return addr, nil
	}

	if mm.constantStrPtr >= CONSTANT_STR_END {
		return -1, fmt.Errorf("string constant memory overflow")
	}

	addr := mm.constantStrPtr

	// Store the string in constant memory
	// The offset should be relative to CONSTANT_START, not CONSTANT_STR_START
	offset := addr - CONSTANT_START
	mm.ConstantMemory[offset] = value

	// Add to map and increment pointer
	mm.constantStringMap[value] = addr
	mm.constantStrPtr++
	return addr, nil
}

func (mm *MemoryManager) Store(address int, value interface{}) error {
	if address < 0 || address >= TOTAL_MEMORY_SIZE {
		return fmt.Errorf("memory access out of bounds: %d", address)
	}

	var segment *[MEMORY_SEGMENT_SIZE]interface{}
	var offset int

	if address >= LOCAL_START && address < LOCAL_START+MEMORY_SEGMENT_SIZE {
		if address < LOCAL_FLOAT_START {
			offset = address - LOCAL_START
		} else {
			offset = address - LOCAL_FLOAT_START + mm.currentSegment.intVarsCount
		}

		if mm.currentSegment == nil {
			return fmt.Errorf("no active function segment")
		}

		if offset >= len(mm.currentSegment.localMemory) {
			return fmt.Errorf("local memory access out of bounds: %d", address)
		}
		// Store directly in the dynamic local memory
		mm.currentSegment.localMemory[offset] = value
		return nil
	}

	switch {
	case address >= CONSTANT_START && address < CONSTANT_START+MEMORY_SEGMENT_SIZE:
		segment = &mm.ConstantMemory
		offset = address - CONSTANT_START
	case address >= TEMP_START && address < TEMP_START+MEMORY_SEGMENT_SIZE:
		segment = &mm.tempMemory
		offset = address - TEMP_START
	case address >= GLOBAL_START && address < GLOBAL_START+MEMORY_SEGMENT_SIZE:
		segment = &mm.globalMemory
		offset = address - GLOBAL_START
	default:
		return fmt.Errorf("invalid memory address: %d", address)
	}

	switch {
	case address >= CONSTANT_STR_START && address <= CONSTANT_STR_END:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("type mismatch: expected string at address %d", address)
		}
	case (address >= GLOBAL_INT_START && address <= GLOBAL_INT_END) ||
		(address >= LOCAL_INT_START && address <= LOCAL_INT_END) ||
		(address >= TEMP_INT_START && address <= TEMP_INT_END) ||
		(address >= CONSTANT_INT_START && address <= CONSTANT_INT_END):

		if floatVal, ok := value.(float64); ok {
			value = int(floatVal)
		}
		if _, ok := value.(int); !ok {
			return fmt.Errorf("type mismatch: expected integer at address %d", address)
		}
	case (address >= GLOBAL_FLOAT_START && address <= GLOBAL_FLOAT_END) ||
		(address >= LOCAL_FLOAT_START && address <= LOCAL_FLOAT_END) ||
		(address >= TEMP_FLOAT_START && address <= TEMP_FLOAT_END) ||
		(address >= CONSTANT_FLOAT_START && address <= CONSTANT_FLOAT_END):

		if intVal, ok := value.(int); ok {
			value = float64(intVal)
		}
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("type mismatch: expected float at address %d", address)
		}
	}
	segment[offset] = value
	return nil
}

func (mm *MemoryManager) Load(address int) (interface{}, error) {
	if address < 0 || address >= TOTAL_MEMORY_SIZE {
		return nil, fmt.Errorf("memory access out of bounds: %d", address)
	}

	var segment *[MEMORY_SEGMENT_SIZE]interface{}
	var offset int

	if address >= LOCAL_START && address < LOCAL_START+MEMORY_SEGMENT_SIZE {
		if address < LOCAL_FLOAT_START {
			offset = address - LOCAL_START
		} else {
			offset = address - LOCAL_FLOAT_START + mm.currentSegment.intVarsCount
		}

		if mm.currentSegment == nil {
			return nil, fmt.Errorf("no active function segment")
		}
		if offset >= len(mm.currentSegment.localMemory) {
			return nil, fmt.Errorf("local memory access out of bounds: %d", address)
		}
		value := mm.currentSegment.localMemory[offset]
		if value == nil {
			return nil, fmt.Errorf("accessing uninitialized memory at address %d", address)
		}
		return value, nil
	}

	switch {
	case address >= CONSTANT_START && address < CONSTANT_START+MEMORY_SEGMENT_SIZE:
		segment = &mm.ConstantMemory
		offset = address - CONSTANT_START
	case address >= TEMP_START && address < TEMP_START+MEMORY_SEGMENT_SIZE:
		segment = &mm.tempMemory
		offset = address - TEMP_START
	case address >= GLOBAL_START && address < MEMORY_SEGMENT_SIZE:
		segment = &mm.globalMemory
		offset = address
	default:
		return nil, fmt.Errorf("invalid memory address: %d", address)
	}

	value := segment[offset]
	if value == nil {
		return nil, fmt.Errorf("accessing uninitialized memory at address %d", address)
	}

	return value, nil
}

func (mm *MemoryManager) PushNewFunctionSegment(isFixed bool, intCount, floatCount int) {
	var size int
	if isFixed {
		size = MEMORY_SEGMENT_SIZE // This is only during function declaration
	} else {
		size = intCount + floatCount
	}

	newSegment := FunctionMemorySegment{
		localMemory:    make([]interface{}, size),
		localIntPtr:    LOCAL_INT_START,
		localFloatPtr:  LOCAL_FLOAT_START,
		intVarsCount:   intCount,
		floatVarsCount: floatCount,
	}

	mm.memoryStack = append(mm.memoryStack, newSegment)
	mm.currentSegment = &mm.memoryStack[len(mm.memoryStack)-1]
}

func (mm *MemoryManager) PopFunctionSegment() error {
	if len(mm.memoryStack) == 0 {
		return fmt.Errorf("no function segments to pop")
	}

	mm.memoryStack = mm.memoryStack[:len(mm.memoryStack)-1]

	if len(mm.memoryStack) > 0 {
		mm.currentSegment = &mm.memoryStack[len(mm.memoryStack)-1]
	} else {
		mm.currentSegment = nil
	}

	return nil
}
