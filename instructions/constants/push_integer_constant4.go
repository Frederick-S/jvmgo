package constants

import (
	"github.com/Frederick-S/jvmgo/instructions/base"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// iconst_4
type PushIntegerConstant4 struct {
	base.NoOperandsInstruction
}

func (pushIntegerConstant4 *PushIntegerConstant4) Execute(frame *runtime_data_area.Frame) {
	frame.GetOperandStack().PushIntegerValue(4)
}
