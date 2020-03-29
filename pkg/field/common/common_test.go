package common

import (
	"testing"

	"github.com/junhwong/goost/pkg/field"
)

func TestREG(t *testing.T) {
	field.Int("message")(12)
}
