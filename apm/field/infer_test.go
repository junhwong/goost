package field

import (
	"fmt"
	"testing"
)

func TestNullPtr(t *testing.T) {
	var x *int
	fmt.Println(InferNumberValue(x))
}
