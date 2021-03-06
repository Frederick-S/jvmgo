package load_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// aload_1
// Load reference from local variable
type ALoad1 struct {
	base_instructions.NoOperandsInstruction
}

func (aLoad1 *ALoad1) Execute(frame *runtime_data_area.Frame) {
	loadReferenceValueAndPush(frame, 1)
}
