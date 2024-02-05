package field

import (
	"fmt"
	"reflect"
	"testing"
)

type cust int

func TestNullPtr(t *testing.T) {
	var x *cust
	fmt.Println(InferPrimitiveValueByReflect(reflect.ValueOf(x)))
	// fmt.Println(InferPrimitiveValueByReflect(reflect.Value{}))
	fmt.Println(InferPrimitiveValue(x))
}

// 2,400,000,000
// 9,007,199,254,740,991
// 9,007,199,254,740,991
// 281,474,976,710,656
