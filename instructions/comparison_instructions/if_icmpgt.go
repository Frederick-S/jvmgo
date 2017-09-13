package comparison_instructions

import (
	"github.com/Frederick-S/jvmgo/instructions/base_instructions"
	"github.com/Frederick-S/jvmgo/runtime_data_area"
)

// if_icmpgt
// Branch if int comparison succeeds.
// Both value1 and value2 must be of type int.
// They are both popped from the operand stack and compared.
// All comparisons are signed. The results of the comparison are as follows:
// if_icmpgt succeeds if and only if value1 > value2
// If the comparison succeeds, the unsigned branchbyte1 and branchbyte2 are used to construct a signed 16-bit offset,
// where the offset is calculated to be (branchbyte1 << 8) | branchbyte2.
// Execution then proceeds at that offset from the address of the opcode of this if_icmp<cond> instruction.
// The target address must be that of an opcode of an instruction within the method that contains this if_icmp<cond> instruction.
// Otherwise, execution proceeds at the address of the instruction following this if_icmp<cond> instruction.
type IfICmpGt struct {
	base_instructions.BranchInstruction
}

func (ifICmpGt *IfICmpGt) Execute(frame *runtime_data_area.Frame) {
	operandStack := frame.GetOperandStack()
	integerValue1 := operandStack.PopIntegerValue()
	integerValue2 := operandStack.PopIntegerValue()

	if integerValue1 > integerValue2 {
		base_instructions.JumpToBranch(frame, ifICmpGt.Offset)
	}
}
