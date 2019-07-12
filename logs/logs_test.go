package logs

import (
	"testing"
)

func TestLog(t *testing.T) {
	root.New().Debug("hello")
	root.New().Info("hello")
}

func TestBuildTemplete(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "${level|-5s} [${time|yyyy-MM-ddTHH:mm:ssS}] ${message}",
		},
		{
			desc: "${level|-5s} [${time|2006-01-02T15:04:05.999Z}] ${message}",
		},
		{
			desc: "${data}",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			buildTemplete(tC.desc)
		})
	}
}
