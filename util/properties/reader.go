package properties

import (
	"bufio"
	"fmt"
	"io"
)

type Reader struct {
	rd   *bufio.Reader
	line int
}

// NewReader 构造一个properties文件读取。
//
// FIXME: 文件编码，优化错误提示
func NewReader(rd io.Reader) *Reader {

	return &Reader{
		rd: bufio.NewReader(rd),
	}
}

func (r *Reader) Next() (key []byte, value []byte, err error) {
	var buf []byte
	var offs int
	r.line++
	buf, _, err = r.rd.ReadLine()
	if err != nil {
		return
	}

	key, offs, err = r.readKey(buf, 0)
	if key == nil && err == nil {
		return r.Next()
	}
	if err != nil {
		return
	}
	value, err = r.readValue(buf, offs)
	return
}
func (r *Reader) readKey(buf []byte, offs int) ([]byte, int, error) {

	a := -1
FOR:
	for offs < len(buf) {
		c := buf[offs]

		switch {
		case a == -1 && isWhiteSpace(c):
			offs++
		case a == -1 && c == '#':
			return nil, 0, nil
		case a == -1 && isKeyPrefix(c):
			a = offs
			offs++
		case a != -1:
			if isKey(c) {
				offs++
				continue FOR
			}
			b := offs
			for offs < len(buf) && isWhiteSpace(buf[offs]) {
				offs++
			}
			if offs < len(buf) && (buf[offs] == ':' || buf[offs] == '=') {
				return buf[a:b], offs + 1, nil
			}
			return nil, 0, fmt.Errorf("Illegal key name: %v", string(buf[offs:]))
		default:
			break FOR
		}
	}
	if a == -1 && offs >= len(buf) {
		return nil, 0, nil // blank line
	}
	return nil, 0, fmt.Errorf("Illegal key name: %v", string(buf[offs:]))
}
func (r *Reader) readValue(buf []byte, offs int) ([]byte, error) {
	a := -1
	b := -1
	inQuotes := false
	for offs < len(buf) && isWhiteSpace(buf[offs]) {
		offs++
	}
	if offs < len(buf) && buf[offs] != '"' {
		a = offs
	}

FOR:
	for offs < len(buf) {
		c := buf[offs]
		switch {
		case a == -1 && c == '"':
			inQuotes = true
			offs++
			a = offs
		case inQuotes:
			for offs < len(buf) {
				c := buf[offs]
				if c == '\\' {
					offs++
					if offs >= len(buf) {
						buf2, _, er := r.rd.ReadLine()
						if er != nil {
							return nil, fmt.Errorf("Value ends unexpectedly")
						}
						buf = append(buf[:offs-1], buf2...)
						r.line++
					} else {
						buf = append(buf[:offs-1], buf[offs:]...)
					}
				} else if c == '"' {
					b = offs
					offs++
					break FOR
				} else {
					offs++
				}
			}
			return nil, fmt.Errorf("Quotation marks are not over")
		case c == '\\':
			offs++
			if offs >= len(buf) {
				buf2, _, er := r.rd.ReadLine()
				if er != nil {
					return nil, fmt.Errorf("Value ends unexpectedly")
				}
				buf = append(buf[:offs-1], buf2...)
				r.line++
			} else {
				buf = append(buf[:offs-1], buf[offs+1:]...)
			}
		case isWhiteSpace(c) || c == '#':
			b = offs
			break FOR
		default:
			offs++
		}
	}
	if a == -1 {
		a = offs
	}
	if b == -1 {
		b = offs
	}

	for offs < len(buf) {
		c := buf[offs]
		if c == '#' {
			break
		} else if isWhiteSpace(c) {
			offs++
		} else {
			return nil, fmt.Errorf("Value ends unexpectedly: %v", string(buf[offs:]))
		}
	}

	return buf[a:b], nil
}
func isWhiteSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\f'
}
func isKeyPrefix(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '$'
}
func isKey(c byte) bool {
	return isKeyPrefix(c) || (c >= '0' && c <= '9') || c == '_' || c == '-' || c == '.'
}
