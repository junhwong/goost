package timestamp

import (
	"fmt"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	nt := Now()
	fmt.Println("CST:", nt.CST().Format(time.RFC3339Nano))
	fmt.Println("UTC:", nt.UTC().Format(time.RFC3339Nano))
	fmt.Println("LOCAL:", nt.Local().Format(time.RFC3339Nano))
	fmt.Println("FORMAT:", nt.Format(time.RFC3339Nano))
}
