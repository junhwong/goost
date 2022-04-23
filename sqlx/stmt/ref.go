package stmt

import (
	"reflect"
	"strings"
)

type Table interface {
	TableName() string
}
type PrimaryKeys interface {
	PrimaryKeyFields() []string
}

////////////
func getNames(obj interface{}, filters []func(string) bool) (map[string]int, reflect.Type, error) {
	typ := reflect.TypeOf(obj)
	if typ.Kind() != reflect.Ptr {
		return nil, nil, newParameterInvalidErr("Must be a ptr")
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		return nil, nil, newParameterInvalidErr("Must be a struct instace")
	}
	names := map[string]int{}
LOOP:
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		jtag := fd.Tag.Get("json")
		if jtag == "-" {
			continue
		}
		name := strings.SplitN(jtag, ",", 2)[0]
		for _, filter := range filters {
			if filter != nil && !filter(name) {
				continue LOOP
			}
		}
		names[name] = i

		// fmt.Println(fd.Name, fd.Type, name)

	}
	return names, typ, nil
}
func getTable(obj interface{}, def string) (string, error) {
	if def != "" {
		return def, nil
	}
	t, _ := obj.(Table)
	if t == nil {
		return "", newParameterInvalidErr("未实现Table接口")
	}
	def = t.TableName()
	if def == "" {
		return "", newParameterInvalidErr("Table名称不能为空")
	}
	return def, nil
}
