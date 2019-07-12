package properties

import (
	"sync"

	"io"

	"github.com/junhwong/goboot/conv"
)

type Properties struct {
	mut  sync.RWMutex
	data map[string]conv.Value
}

var notfound = conv.Warp(nil)

func (props *Properties) Set(key string, value interface{}) conv.Value {
	props.mut.Lock()
	v := conv.Warp(value)
	props.data[key] = v
	props.mut.Unlock()
	return v
}

func (props *Properties) Get(key string) conv.Value {
	props.mut.RLock()
	v, ok := props.data[key]
	props.data[key] = v
	props.mut.RUnlock()
	if !ok {
		return notfound
	}
	return v
}

func (props *Properties) Load(r io.Reader, prefix ...string) error {
	props.mut.RLock()
	defer props.mut.RUnlock()
	rd := NewReader(r)
	var p string
	if len(prefix) > 0 && prefix[0] != "" {
		p = prefix[0]
	}
	for {
		k, v, err := rd.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		key := string(k)
		value := string(v) // FIXME: 编码解码
		if len(key) == 0 || len(value) == 0 {
			continue
		}
		if p != "" {
			key = p + "." + key
		}
		props.data[key] = conv.Warp(value)
	}
	return nil
}
