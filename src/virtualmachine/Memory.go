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
	globalMemory   []interface{}
	GlobalIntPtr   int
	GlobalFloatPtr int

	tempMemory   []interface{}
	TempIntPtr   int
	TempFloatPtr int

	constantIntPtr   int
	constantFloatPtr int
	constantStrPtr   int

	// map for constant reusing / not restoring the same constant
	ConstantMapLoad  map[int]interface{}
	ConstantMapStore map[interface{}]int

	memoryStack    []FunctionMemorySegment
	currentSegment *FunctionMemorySegment
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		GlobalIntPtr:     GLOBAL_INT_START,
		GlobalFloatPtr:   GLOBAL_FLOAT_START,
		TempIntPtr:       TEMP_INT_START,
		TempFloatPtr:     TEMP_FLOAT_START,
		constantIntPtr:   CONSTANT_INT_START,
		constantFloatPtr: CONSTANT_FLOAT_START,
		constantStrPtr:   CONSTANT_STR_START,
		ConstantMapLoad:  make(map[int]interface{}),
		ConstantMapStore: make(map[interface{}]int),
		memoryStack:      make([]FunctionMemorySegment, 0),
		currentSegment:   nil,
	}
}

func (mm *MemoryManager) InitializeMemory() {
	globalFloat := mm.GlobalFloatPtr
	globalInt := mm.GlobalIntPtr
	globalSize := globalFloat + globalInt

	tempInt := mm.TempIntPtr
	tempFloat := mm.TempFloatPtr
	tempSize := tempInt + tempFloat

	mm.globalMemory = make([]interface{}, globalSize)
	mm.tempMemory = make([]interface{}, tempSize)
}

func (mm *MemoryManager) AllocateGlobal(varType shared.Type) (int, error) {
	switch varType {
	case shared.TypeInt:
		if mm.GlobalIntPtr >= GLOBAL_INT_END {
			return -1, fmt.Errorf("global integer memory overflow")
		}
		addr := mm.GlobalIntPtr
		mm.GlobalIntPtr++
		return addr, nil
	case shared.TypeFloat:
		if mm.GlobalFloatPtr >= GLOBAL_FLOAT_END {
			return -1, fmt.Errorf("global float memory overflow")
		}
		addr := mm.GlobalFloatPtr
		mm.GlobalFloatPtr++
		return addr, nil
	default:
		return -1, fmt.Errorf("unsupported type for global allocation")
	}
}

func (mm *MemoryManager) AllocateTemp(varType shared.Type) (int, error) {
	switch varType {
	case shared.TypeInt:
		if mm.TempIntPtr >= TEMP_INT_END {
			return -1, fmt.Errorf("temporary integer memory overflow")
		}
		addr := mm.TempIntPtr
		mm.TempIntPtr++
		return addr, nil
	case shared.TypeFloat:
		if mm.TempFloatPtr >= TEMP_FLOAT_END {
			return -1, fmt.Errorf("temporary float memory overflow")
		}
		addr := mm.TempFloatPtr
		mm.TempFloatPtr++
		return addr, nil
	default:
		return -1, fmt.Errorf("unsupported type for temporary allocation")
	}
}

func (mm *MemoryManager) AllocateConstant(value string) (int, error) {
	if addr, exists := mm.ConstantMapStore[value]; exists {
		return addr, nil
	}

	if intVal, ok := strconv.Atoi(value); ok == nil {
		if mm.constantIntPtr >= CONSTANT_INT_END {
			return -1, fmt.Errorf("constant integer memory overflow")
		}
		addr := mm.constantIntPtr
		// fmt.Println("Allocating value ", intVal, "at ", addr-CONSTANT_START)
		mm.ConstantMapLoad[addr-CONSTANT_START] = intVal
		mm.ConstantMapStore[value] = addr
		mm.constantIntPtr++
		return addr, nil
	}

	if floatVal, ok := strconv.ParseFloat(value, 64); ok == nil {
		if mm.constantFloatPtr >= CONSTANT_FLOAT_END {
			return -1, fmt.Errorf("constant integer memory overflow")
		}
		addr := mm.constantFloatPtr
		mm.ConstantMapLoad[addr-CONSTANT_START] = floatVal
		mm.ConstantMapStore[value] = addr
		mm.constantFloatPtr++
		return addr, nil
	}

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

func (mm *MemoryManager) AllocateStringAddress(value string) (int, error) {
	if addr, exists := mm.ConstantMapStore[value]; exists {
		return addr, nil
	}

	if mm.constantStrPtr >= CONSTANT_STR_END {
		return -1, fmt.Errorf("string constant memory overflow")
	}

	addr := mm.constantStrPtr

	// Store the string in constant memory
	offset := addr - CONSTANT_START
	mm.ConstantMapLoad[offset] = value

	// Add to map and increment pointer
	mm.ConstantMapStore[value] = addr
	mm.constantStrPtr++
	return addr, nil
}

func (mm *MemoryManager) Store(address int, value interface{}) error {
	if address < 0 || address >= TOTAL_MEMORY_SIZE {
		return fmt.Errorf("memory access out of bounds: %d", address)
	}
	// fmt.Println("entered here:  with address ", address, "and ", value)

	var segment *[]interface{}
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
	case address >= TEMP_START && address < TEMP_START+MEMORY_SEGMENT_SIZE:
		segment = &mm.tempMemory
		if address >= TEMP_INT_START && address < TEMP_FLOAT_START {
			offset = address - TEMP_START
		} else if address >= TEMP_FLOAT_START {
			offset = address - TEMP_FLOAT_START + mm.TempIntPtr
		}
	case address >= GLOBAL_START && address < GLOBAL_START+MEMORY_SEGMENT_SIZE:
		segment = &mm.globalMemory
		if address >= GLOBAL_INT_START && address < GLOBAL_FLOAT_START {
			offset = address - GLOBAL_START
		} else if address >= GLOBAL_FLOAT_START {
			offset = address - GLOBAL_FLOAT_START + mm.GlobalIntPtr
		}
	default:
		return fmt.Errorf("invalid memory address: %d", address)
	}

	(*segment)[offset] = value
	return nil
}

func (mm *MemoryManager) Load(address int) (interface{}, error) {
	if address < 0 || address >= TOTAL_MEMORY_SIZE {
		return nil, fmt.Errorf("memory access out of bounds: %d", address)
	}

	var segment *[]interface{}
	var offset int

	if address >= LOCAL_START && address < LOCAL_START+MEMORY_SEGMENT_SIZE {
		if address < LOCAL_FLOAT_START {
			offset = address - LOCAL_START
		} else {
			offset = (address - LOCAL_FLOAT_START) + mm.currentSegment.intVarsCount
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
		offset := address - CONSTANT_START
		if _, exists := mm.ConstantMapLoad[offset]; !exists {
			return nil, fmt.Errorf("trying to retrieve from uninitialized memory address")
		}
		return mm.ConstantMapLoad[offset], nil
	case address >= TEMP_START && address < TEMP_START+MEMORY_SEGMENT_SIZE:
		segment = &mm.tempMemory
		if address >= TEMP_INT_START && address < TEMP_FLOAT_START {
			offset = address - TEMP_START
		} else if address >= TEMP_FLOAT_START {
			offset = address - TEMP_FLOAT_START + mm.TempIntPtr
		}
	case address >= GLOBAL_START && address < MEMORY_SEGMENT_SIZE:
		segment = &mm.globalMemory
		if address >= GLOBAL_INT_START && address < GLOBAL_FLOAT_START {
			offset = address - GLOBAL_START
		} else if address >= GLOBAL_FLOAT_START {
			offset = address - GLOBAL_FLOAT_START + mm.GlobalIntPtr
		}
	default:
		return nil, fmt.Errorf("invalid memory address: %d", address)
	}

	value := (*segment)[offset]
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
