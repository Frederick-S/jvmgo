package runtime_data_area

import "github.com/Frederick-S/jvmgo/runtime_data_area/heap"

type Variable struct {
	value          int32
	referenceValue *heap.Object
}
