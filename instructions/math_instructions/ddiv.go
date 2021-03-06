package math_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// ddiv
// Divide double
type DDiv struct {
	base_instructions.NoOperandsInstruction
}

func (dDiv *DDiv) Execute(frame *runtime_data_area.Frame) {
	operandStack := frame.GetOperandStack()
	doubleValue2 := operandStack.PopDoubleValue()
	doubleValue1 := operandStack.PopDoubleValue()

	operandStack.PushDoubleValue(doubleValue1 / doubleValue2)
}
