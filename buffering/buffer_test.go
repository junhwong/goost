package buffering

import (
	"fmt"
	"testing"
)

func TestGroup(t *testing.T) {
	var buf Buffer
	buf.WriteByte('0')
	buf.WriteByte('1')
	buf.WriteByte('2')
	fmt.Printf("buf.buf: %s\n", buf.buf)
	buf.Write([]byte{'3', '4', '5'})
	buf.WriteByte('6')
	fmt.Printf("buf.buf: %s\n", buf.buf)

}
