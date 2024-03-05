package field

import (
	"fmt"
	"testing"
)

type cust int

func TestNullPtr(t *testing.T) {
	var x *cust
	// fmt.Println(InferPrimitiveValueByReflect(reflect.ValueOf(x)))
	// fmt.Println(InferPrimitiveValueByReflect(reflect.Value{}))
	fmt.Println(Any("", x))

	// var y *int
	// // fmt.Println(InferPrimitiveValueByReflect(reflect.ValueOf(y)))
	// // fmt.Println(InferPrimitiveValueByReflect(reflect.Value{}))
	// fmt.Println(InferPrimitiveValue(y))
}

// 2,400,000,000
// 9,007,199,254,740,991
// 9,007,199,254,740,991
// 281,474,976,710,656
