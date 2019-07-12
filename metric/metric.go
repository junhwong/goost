package metric

type Metric interface {
	GetName() string
	GetValue() interface{}
}

type StringMetric struct {
	Name  string
	Value string
}

func (sm *StringMetric) GetName() string {
	return sm.Name
}
func (sm *StringMetric) GetValue() interface{} {
	return sm.Value
}
func String(name, val string) Metric {
	return &StringMetric{name, val}
}
