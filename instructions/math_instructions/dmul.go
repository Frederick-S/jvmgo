package math_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// dmul
// Multiply double
type DMul struct {
	base_instructions.NoOperandsInstruction
}

func (dMul *DMul) Execute(frame *runtime_data_area.Frame) {
	operandStack := frame.GetOperandStack()
	doubleValue2 := operandStack.PopDoubleValue()
	doubleValue1 := operandStack.PopDoubleValue()

	operandStack.PushDoubleValue(doubleValue1 * doubleValue2)
}
