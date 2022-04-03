package json

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type jsonParser struct {
	reader *bufio.Reader
}

func (r *jsonParser) readValue() (interface{}, error) {
	for {
		c, err := r.reader.ReadByte()
		if err != nil {
			return "", err
		}
		switch c {
		case '\\':
		case '"':

		default:
			return "", fmt.Errorf("非法字符")
		}
	}
}
func (r *jsonParser) readString() (string, error) {
	for {
		c, err := r.reader.ReadByte()
		if err != nil {
			return "", err
		}
		switch c {
		case '\\':
		case '"':

		default:
			return "", fmt.Errorf("非法字符")
		}
	}
}
func (r *jsonParser) readStruct() (map[string]interface{}, error) {
	dst := map[string]interface{}{}
	for {
		c, err := r.reader.ReadByte()
		if err != nil {
			return nil, err
		}
		switch c {
		case '}':
			break
		case '"':
			key, _ := r.readString()
			// read :
			val, _ := r.readValue()
			dst[key] = val
		case ',':
		default:
			return nil, fmt.Errorf("非法字符")
		}
	}
}
func (r *jsonParser) ToJson() {

	sr := strings.NewReader(`{
		"username":"test",
		"password":"123456",
		"client_id":"1"
	}`)
	r.reader = bufio.NewReader(sr)
	for {
		c, err := r.reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		switch c {
		case '{':
			r.reader.UnreadByte()
			r.reader.Peek(10)
		}
	}

}
