package buffering

import (
	"fmt"
	"io"
)

type Buffer struct {
	buf   []byte
	start int
	end   int
}

func (b *Buffer) SetDirect(buf []byte) {
	PutBytes(b.buf)
	b.buf = buf
	b.start = 0
	b.end = len(buf)
}

func (b *Buffer) grow(n int) {
	if n < 0 {
		panic("negative grow")
	}
	tot := cap(b.buf)
	if n+b.end <= tot {
		return
	}

	if n < 2048 { // 提升效率
		n = 2048
	}
	buf := make([]byte, n+b.end)
	if b.end != 0 {
		b.end = copy(buf, b.buf[b.start:b.end])
	}
	b.start = 0
	b.buf = buf
}
func (b *Buffer) WriteByte(c byte) error {
	b.grow(1)
	b.buf[b.end] = c
	b.end++
	return nil
}
func (b *Buffer) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	b.grow(len(p))
	n := copy(b.buf[b.end:], p)
	b.end += n
	return n, nil
}
func (b *Buffer) WriteString(s string) (int, error) {
	b.grow(len(s))
	n := copy(b.buf[b.end:], s)
	b.end += n
	return n, nil
}
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	// buf := b.buf[b.start:b.end]
	// n, err := w.Write(buf)
	// b.start += n
	// cnt := int64(n)
	// return cnt, err
	panic("not implemented")
}
func (b *Buffer) ReadFromN(r io.Reader, n int) (written int64, err error) {
	var remaining int
	if n > 0 {
		remaining = n
		b.grow(n)
	} else {
		remaining = 1
		b.grow(512)
	}

	for remaining > 0 {
		buf := b.buf[b.end:]
		c, err := r.Read(buf)
		written += int64(c)
		b.end += c
		if err == io.EOF {
			return written, nil
		}
		if err != nil {
			return written, err
		}
		if n <= 0 {
			b.grow(512)
		} else {
			remaining -= c
		}
	}
	return
}
func (b *Buffer) empty() bool {
	return b.end == b.start
}
func (b *Buffer) Offset() int {
	return b.start
}
func (b *Buffer) Len() int {
	return b.end - b.start
}

func (b *Buffer) Truncate(n int) {
	if n < 0 {
		panic("Buffer.Truncate: negative count")
	}
	if n > b.Len() {
		panic("Buffer.Truncate: count > length")
	}
	b.end -= n
}

func (b *Buffer) Next(n int) []byte {
	if n < 0 {
		panic("Buffer.Next: negative count")
	}
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.start : b.start+n]
	b.start += n
	return data
}
func (b *Buffer) Prev(n int) []byte {
	if n < 0 {
		panic("Buffer.Prev: negative count")
	}
	if b.start-n < 0 {
		fmt.Printf("b.start: %v\n", b.start)
		fmt.Printf("n: %v\n", n)
		panic("Buffer.Prev: negative count")
	}
	data := b.buf[b.start-n : b.start]
	b.start -= n
	return data
}
func (b *Buffer) Reset() {
	b.start = 0
	b.end = 0
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		return 0, nil
	}
	if n > b.Len() {
		n = b.Len()
	}
	if n == 0 {
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.start:b.start+n])
	b.start += n
	return
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		return 0, io.EOF
	}
	c := b.buf[b.start]
	b.start++
	return c, nil
}

func (b *Buffer) Bytes() []byte {
	return b.buf[b.start:b.end]
}

func (b *Buffer) String() string {
	return string(b.Bytes())
}
