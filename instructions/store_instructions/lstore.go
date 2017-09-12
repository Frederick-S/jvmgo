package store_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// lstore
// Store long into local variable.
// The index is an unsigned byte.
// Both index and index+1 must be indices into the local variable array of the current frame.
// The value on the top of the operand stack must be of type long.
// It is popped from the operand stack, and the local variables at index and index+1 are set to value.
type LStore struct {
	base_instructions.Index8Instruction
}

func (lStore *LStore) Execute(frame *runtime_data_area.Frame) {
	popLongValueAndStore(frame, uint(lStore.Index))
}

func popLongValueAndStore(frame *runtime_data_area.Frame, index uint) {
	value := frame.GetOperandStack().PopLongValue()
	frame.GetLocalVariables().SetLongValue(index, value)
}
