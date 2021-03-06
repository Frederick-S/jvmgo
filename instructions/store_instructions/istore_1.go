package store_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// istore_1
// Store int into local variable
type IStore1 struct {
	base_instructions.NoOperandsInstruction
}

func (iStore1 *IStore1) Execute(frame *runtime_data_area.Frame) {
	popIntegerValueAndStore(frame, 1)
}
