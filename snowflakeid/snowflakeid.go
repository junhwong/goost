package snowflakeid

import (
	"context"
	"fmt"
)

var (
	defaultGen IDGen
)

func init() {
	gen, err := DefaultBuilder.Build()
	if err != nil {
		panic(err)
	}
	_ = SetDefault(gen)
}

type IDGen func(context.Context) int64

// GenerateID 返回一个新的ID
func GenerateID(ctx context.Context) int64 {
	return defaultGen(ctx)
}

func SetDefault(gen IDGen) error {
	if gen == nil {
		return fmt.Errorf("gen cannot be nil")
	}
	defaultGen = gen
	return nil
}

var DefaultBuilder = Builder{
	Layout:    DefaultLayout,
	WorkerGen: WorkerWithPodIDOrHostname,
	TimeGen:   UTCMillisecond,
}
