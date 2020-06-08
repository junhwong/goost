package tcc

type State struct {
	Name             string                 `json:"n"`
	Phase            string                 `json:"p"`
	TransactionID    uint64                 `json:"tid"`
	TransactionState uint64                 `json:"ts"`
	Entries          map[string]interface{} `json:"es"`
}

func (st *State) AddEntry(key string, val interface{}) {
	st.Entries[key] = val
}

func NewState(ctx TransactionContext) State {
	return State{
		Name:             ctx.Name(),
		Phase:            ctx.Phase(),
		TransactionID:    ctx.TransactionID(),
		TransactionState: ctx.TransactionState(),
		Entries:          make(map[string]interface{}),
	}
}
