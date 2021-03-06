package load_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// iload_2
// Load int from local variable
type ILoad2 struct {
	base_instructions.NoOperandsInstruction
}

func (iLoad2 *ILoad2) Execute(frame *runtime_data_area.Frame) {
	loadIntegerValueAndPush(frame, 2)
}
