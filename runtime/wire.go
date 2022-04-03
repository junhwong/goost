package runtime

import (
	"container/list"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	_errType    = reflect.TypeOf((*error)(nil)).Elem()
	_inPtrType  = reflect.TypeOf((*In)(nil))
	_inType     = reflect.TypeOf(In{})
	_outPtrType = reflect.TypeOf((*Out)(nil))
	_outType    = reflect.TypeOf(Out{})
)

type nodeKind int

const (
	nodeFunc   nodeKind = iota // 函数
	nodeResult                 // 包含依赖值得结果
	nodeParam                  // 函数参数
	nodeField                  // 结构字段
)

// 标签名称
var TagName = "inject"

func isError(t reflect.Type) bool {
	return t.Implements(_errType)
}

// Returns true if t embeds e or if any of the types embedded by t embed e.
func embedsType(i interface{}, e reflect.Type) bool {
	// TODO: this function doesn't consider e being a pointer.
	// given `type A foo { *In }`, this function would return false for
	// embedding dig.In, which makes for some extra error checking in places
	// that call this funciton. Might be worthwhile to consider reflect.Indirect
	// usage to clean up the callers.

	if i == nil {
		return false
	}

	// maybe it's already a reflect.Type
	t, ok := i.(reflect.Type)
	if !ok {
		// take the type if it's not
		t = reflect.TypeOf(i)
	}

	// We are going to do a breadth-first search of all embedded fields.
	types := list.New()
	types.PushBack(t)
	for types.Len() > 0 {
		t := types.Remove(types.Front()).(reflect.Type)

		if t == e {
			return true
		}

		if t.Kind() != reflect.Struct {
			continue
		}

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Anonymous {
				types.PushBack(f.Type)
			}
		}
	}

	// If perf is an issue, we can cache known In objects and Out objects in a
	// map[reflect.Type]struct{}.
	return false
}

type container struct {
	typedIndex  map[reflect.Type]*node
	namedIndex  map[string]*node
	unbindNodes []*node
	nodes       []*node
	callNodes   []*node
	err         error
}

func (c *container) newNode(kind nodeKind, ctype reflect.Type) (n *node) {

	n = &node{
		ctype: ctype,
		kind:  kind, // input
	}
	c.nodes = append(c.nodes, n)
	return n
}
func (c *container) mkerr(err error) error {
	c.err = err
	return c.err
}

// 安装一个函数
func (c *container) setup(constructor interface{}, calldept int) error {
	if c.err != nil {
		return c.err
	}
	ctype := reflect.TypeOf(constructor)
	if ctype == nil {
		return c.mkerr(errors.New("can't provide an untyped nil"))
	}
	if ctype.Kind() != reflect.Func {
		return c.mkerr(fmt.Errorf("must provide constructor function, got %v (type %v)", constructor, ctype))
	}

	n := c.newNode(nodeFunc, ctype)
	n.cotr = constructor
	c.callNodes = append(c.callNodes, n)

	i := ctype.NumOut() - 1
	var rtype reflect.Type
	for i > -1 {
		t := ctype.Out(i)
		i--

		if isError(t) {
			if n.checkerr {
				return c.mkerr(errors.New("返回结果只能有个一个错误类型且必须在最后"))
			}
			n.checkerr = true
			continue
		}
		if rtype != nil {
			return c.mkerr(errors.New("返回结果只能有一个非错误类型"))
		}
		rtype = t // 该类型
	}

	if rtype != nil {
		// TODO: 检查相同类型只能出现一次，否则只能命名提供
		n.result = c.reg(n.name, rtype)
		n.result.dep = n
		//c.checkInject(n.returnNode)
		//check inject
	}

	// n.params = []*node{}
	for i := 0; i < ctype.NumIn(); i++ {
		// 2 parm
		p := c.newNode(nodeParam, ctype.In(i))
		p.index = i
		// ctype.IsVariadic()
		if err := p.checkParam(); err != nil {
			return c.mkerr(err)
		}
		n.params = append(n.params, p)
		c.unbindNodes = append(c.unbindNodes, p)

	}
	return nil
}

func (c *container) checkInject(n *node) {
	ctype := n.ctype
	if ctype.Kind() == reflect.Ptr {
		ctype = ctype.Elem()
	}
	if ctype.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < ctype.NumField(); i++ {
		f := ctype.Field(i)
		tag := f.Tag.Get(TagName)
		if tag == "" {
			continue
		}

		cf := c.newNode(nodeField, f.Type)
		cf.index = i
		for i, it := range strings.Split(tag, ",") {
			it = strings.TrimSpace(it)
			if i == 0 {
				cf.name = it
				continue
			}
			switch it {
			case "required":
				cf.required = true
			case "group":
				//TODO: check

			}
		}
		n.fields = append(n.fields, cf)
		// dep := c.getProv(cf.name, cf.ctype)
		// if dep != nil {
		// 	cf.deps = append(cf.deps, dep)
		// } else {
		// 	c.undep = append(c.undep, cf)
		// }
		// reflect.ValueOf(nil).Field(0).Set(reflect.ValueOf(nil))
		fmt.Println("f.Tag", tag)
	}
	fmt.Println("checkInject: Struct")
}

func (c *container) reg(name string, ctype reflect.Type) (v *node) {
	v = c.namedIndex[name]
	if v == nil {
		v = c.typedIndex[ctype]
	}
	// else if c.values[ctype]!=nil {
	// 	panic("已经注册")
	// }
	if v != nil {
		panic("已经注册")
	}
	v = c.newNode(nodeResult, ctype) // 函数执行结果,也将作为值提供

	if name != "" {
		delete(c.typedIndex, ctype)
		c.namedIndex[name] = v
	} else {
		c.typedIndex[ctype] = v
	}
	v.name = name
	return v
}

// func (c *container) getVal(ctype reflect.Type) interface{} {
// 	v := c.typedIndex[ctype]
// 	if v == nil {
// 		return nil
// 	}
// 	return v.returnNode
// }
func (c *container) getProv(name string, ctype reflect.Type) (v *node) {
	v = c.namedIndex[name]
	if v != nil {
		return
	}
	v = c.typedIndex[ctype]
	return
}
func (c *container) bind() {
	if c.err != nil {
		return
	}
	news := []*node{}
	for _, it := range c.unbindNodes {
		n := c.getProv(it.name, it.ctype)
		if n == nil {
			news = append(news, it)
			continue
		}
		it.dep = n
	}
	c.unbindNodes = news
}

func (c *container) call() error {
	if c.err != nil {
		return c.err
	}
	if len(c.unbindNodes) != 0 {
		return c.mkerr(fmt.Errorf("有未绑定的依赖节点"))
	}
	for _, it := range c.callNodes {
		fmt.Println("call type:", it.ctype)
		if err := it.call(c); err != nil {
			return err
		}
	}
	return nil
}

// 节点
type node struct {
	cotr      interface{}  // 函数或值
	ctype     reflect.Type // 类型
	kind      nodeKind     // 节点类型
	params    []*node      // 函数参数
	checkerr  bool         // 是否检查返回错误(针对函数)
	name      string       // 绑定名称
	result    *node        // 函数返回值
	index     int          // 参数或字段索引
	fields    []*node      // 结构体字段
	required  bool         // 依赖是否允许为空
	grouped   bool         // 是否是组
	groupDeps []*node      // 组结果依赖
	dep       *node        // 依赖项
	provided  bool         // 是否已经处理
}

func (n *node) checkParam() error {
	//TODO: 参数依赖
	// 1 参数不能有简单类型
	// 2 结构参数构造
	switch n.ctype.Kind() {
	case reflect.Bool, reflect.String, reflect.Int:
		return fmt.Errorf("参数不能是简单类型: %v", n.ctype)
	case reflect.Struct:

	}
	return nil
}
func (n *node) provide(c *container) error {
	if n.provided {
		return nil
	}
	if n.kind != nodeResult {
		return fmt.Errorf("不是值节点")
	}

	if n.dep == nil {
		return fmt.Errorf("未绑定依赖函数: %v", n.ctype)
	}
	err := n.dep.call(c)
	n.provided = true // 有什么用？
	return err
}
func (n *node) call(c *container) error {
	if n.provided {
		return nil
	}
	if len(c.unbindNodes) > 0 {
		return fmt.Errorf("还要未绑定的依赖")
	}
	if n.kind != nodeFunc {
		return fmt.Errorf("不是函数节点")
	}
	params := []reflect.Value{}
	for _, it := range n.params {
		if err := it.dep.provide(c); err != nil {
			return err
		}
		v := it.dep.cotr
		params = append(params, reflect.ValueOf(v))
	}
	n.provided = true
	cotr := reflect.ValueOf(n.cotr)
	returns := cotr.Call(params)
	var rerr interface{}
	var val interface{}
	if n.checkerr && n.result != nil {
		rerr = returns[1].Interface()
		val = returns[0].Interface()
	} else if n.result != nil {
		val = returns[0].Interface()
	} else if n.checkerr {
		rerr = returns[0].Interface()
	}
	if rerr != nil {
		return fmt.Errorf("返回错误: %v", rerr)
	}
	if val != nil && n.result != nil {
		n.result.cotr = val
		// TODO: post hook
	}
	return nil
}
