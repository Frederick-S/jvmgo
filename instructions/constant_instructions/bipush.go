package constant_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// bipush
// Push byte
type BIPush struct {
	value int8
}

func (bIPush *BIPush) FetchOperands(bytecodeReader *base_instructions.BytecodeReader) {
	bIPush.value = bytecodeReader.ReadInt8()
}

func (bIPush *BIPush) Execute(frame *runtime_data_area.Frame) {
	integer := int32(bIPush.value)

	frame.GetOperandStack().PushIntegerValue(integer)
}
