package Storer

import (
	"encoding/gob"
	"fmt"
	"os"
	"pogo/src/semantic"
	"pogo/src/shared"
	"pogo/src/virtualmachine"
)

type SerializedVMData struct {
	Quadruples     []shared.Quadruple
	Functions      map[string]shared.FunctionInfo
	ConstantMemory [4000]interface{}
}

func SaveCompiledData(quads []shared.Quadruple, SymbolTable *semantic.SymbolTable, memoryManager *virtualmachine.MemoryManager, filename string) error {
	functions := make(map[string]shared.FunctionInfo)

	for name, symbol := range SymbolTable.GetGlobalScope() {
		if function, ok := symbol.(shared.Function); ok {
			functions[name] = shared.FunctionInfo{
				Name:           function.Name,
				StartQuad:      function.StartQuad,
				IntVarsCount:   function.IntVarsCounter,
				FloatVarsCount: function.FloatVarsCounter,
				Parameters:     function.Parameters,
			}
		}
	}

	vmData := SerializedVMData{
		Quadruples:     quads,
		Functions:      functions,
		ConstantMemory: memoryManager.ConstantMemory,
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(vmData); err != nil {
		return fmt.Errorf("error encoding data: %v", err)
	}

	return nil
}

func LoadCompiledData(filename string) (*virtualmachine.VirtualMachine, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	var vmData SerializedVMData
	if err := decoder.Decode(&vmData); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	memManager := virtualmachine.NewMemoryManager()
	memManager.ConstantMemory = vmData.ConstantMemory
	vm := virtualmachine.NewVirtualMachine(vmData.Quadruples, memManager)

	vm.Functions = vmData.Functions

	return vm, nil
}
