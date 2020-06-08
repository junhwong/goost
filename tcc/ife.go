package tcc

type TransactionContext interface {
	Name() string
	Phase() string
	TransactionID() uint64
	TransactionState() uint64 // try|confirm|cancel
	Entries() interface{}
}

type TryFn func(ctx TransactionContext) (state State, err error)

type Provider interface {
	Name() string
}

func Try(ctx TransactionContext, phase string, fn TryFn) error {

	state, err := fn(ctx)
	if err != nil {
		return err
	}
	if state.Name == "" {

	}
	return nil
}
