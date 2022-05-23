package benchmarks

import (
	"testing"

	"github.com/junhwong/goost/apm"
)

func BenchmarkAccumulatedContext(b *testing.B) {
	b.Logf("Logging with some accumulated context.")
	b.Run("goost/apm", func(b *testing.B) {
		// logger := newZapLogger(zap.DebugLevel).With(fakeFields()...)
		logger := apm.Default()
		//
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		b.Cleanup(apm.Done)
	})
}
