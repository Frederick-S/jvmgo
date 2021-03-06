package comparison_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// fcmpl
// Compare float
type FCmpl struct {
	base_instructions.NoOperandsInstruction
}

func (fCmpl *FCmpl) Execute(frame *runtime_data_area.Frame) {
	operandStack := frame.GetOperandStack()
	floatValue2 := operandStack.PopFloatValue()
	floatValue1 := operandStack.PopFloatValue()

	if floatValue1 > floatValue2 {
		operandStack.PushIntegerValue(1)
	} else if floatValue1 == floatValue2 {
		operandStack.PushIntegerValue(0)
	} else if floatValue1 < floatValue2 {
		operandStack.PushIntegerValue(-1)
	} else {
		operandStack.PushIntegerValue(-1)
	}
}
