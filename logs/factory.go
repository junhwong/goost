package logs

// var root = &Factory{
// 	Handler: NewFormatHandler(DEBUG, "${level|-5s} [${time|yyyy-MM-ddTHH:mm:ssS}] ${message}\n", nil),
// }

// // Factory
// type Factory struct {
// 	Prefix  string
// 	Level   Level
// 	Tags    Fields
// 	Data    Fields
// 	Handler Handler
// 	Clone   func() *Factory
// }

// func (fa *Factory) New() *Entry {
// 	entry := &Entry{
// 		Prefix:  fa.Prefix,
// 		Handler: fa,
// 		Level:   fa.Level,
// 		Tags:    make(Fields),
// 		Data:    make(Fields),
// 		Fields:  make([]*Field, 0),
// 	}
// 	for k, v := range fa.Tags {
// 		entry.Tags[k] = v
// 	}
// 	for k, v := range fa.Data {
// 		entry.Data[k] = v
// 	}
// 	return entry
// }

// func (fa *Factory) Handle(entry *Entry) {
// 	fa.Handler.Handle(entry)
// }

// // ####### 实现 Logger 接口方法
// func clone(fa *Factory) *Factory {
// 	if fa.Clone != nil {
// 		return fa.Clone()
// 	}
// 	return &Factory{
// 		Prefix: fa.Prefix,
// 		Level:  fa.Level,
// 	}
// }
// func (fa *Factory) WithPrefix(prefix string, joinParent ...bool) *Factory {
// 	if len(joinParent) > 0 && joinParent[0] && fa.Prefix != "" && prefix != "" {
// 		prefix = fa.Prefix + "." + prefix
// 	}
// 	c := clone(fa)
// 	c.Prefix = prefix
// 	return c
// }
// func (fa *Factory) WithLevel(lvl Level) *Factory {
// 	c := clone(fa)
// 	c.Level = lvl
// 	return c
// }
// func (fa *Factory) WithContext(ctx context.Context) *Factory {
// 	c := clone(fa)
// 	// TODO:
// 	return c
// }
// func (fa *Factory) WithFields(fields *Field, replace ...bool) *Factory {
// 	c := clone(fa)
// 	// TODO:
// 	return c
// }
