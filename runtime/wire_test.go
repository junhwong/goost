package runtime

import (
	"fmt"
	"reflect"
	"testing"
)

type demosettings struct {
	In
}

type insImpl struct {
	s string
	// DB interface{} `inject:"abc,required,group"`
}

func newIns(s demosettings) *insImpl {
	return &insImpl{s: "s"}
}

func run(i *insImpl) {
	fmt.Println("===执行结果: ", i.s)
}

func TestProvide(t *testing.T) {
	c := &container{
		typedIndex: make(map[reflect.Type]*node),
		namedIndex: make(map[string]*node),
	}
	c.setup(newIns, 0)
	c.setup(run, 0)
	c.bind()
	// for _, dep := range c.unbindNodes {
	// 	t.Log("undep:", dep.ctype.Kind(), dep.name)
	// }
	t.Log("end call:", c.call())

}

// type cvalue struct {
// 	val      interface{}
// 	name     string
// 	ctype    reflect.Type
// 	grouped  bool
// 	provided bool
// 	fields   []*cfield
// }

//in option: type;name;nullable;
//pi option: type;interfaces;name;nullable;
