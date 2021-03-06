package math_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// fadd
// Add float
type FAdd struct {
	base_instructions.NoOperandsInstruction
}

func (fAdd *FAdd) Execute(frame *runtime_data_area.Frame) {
	operandStack := frame.GetOperandStack()
	floatValue2 := operandStack.PopFloatValue()
	floatValue1 := operandStack.PopFloatValue()

	operandStack.PushFloatValue(floatValue1 + floatValue2)
}
