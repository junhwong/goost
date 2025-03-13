package field

import (
	"errors"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

var unmarshaler proto.UnmarshalOptions
var marshaler proto.MarshalOptions
var bits = 4

func MarshalProto(f *Field) ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	b1 := make([]byte, bits)
	b2, err := marshaler.MarshalAppend(b1, f.Schema)
	if err != nil {
		return nil, err
	}
	n := len(b2)
	setLen32(b2, n-bits, 0)
	b2 = protowire.AppendFixed32(b2, 0)
	// b2 = b1[:n+bits]
	b2, err = marshaler.MarshalAppend(b2, f.Value)
	if err != nil {
		return nil, err
	}
	setLen32(b2, len(b2)-bits-n, n)
	return b2, nil
}
func UnmarshalProto(b2 []byte, f *Field) error {
	if f == nil {
		return errors.New("f cannot be nil")
	}
	v, n := protowire.ConsumeFixed32(b2)
	if n != bits {
		return errors.New("invalid length")
	}
	n = int(v)
	b2 = b2[bits:]

	err := unmarshaler.Unmarshal(b2[:n], f.Schema)
	if err != nil {
		return err
	}
	b2 = b2[n:]
	v, n = protowire.ConsumeFixed32(b2)
	if n != bits {
		return errors.New("invalid length")
	}
	n = int(v)
	b2 = b2[bits:]
	err = unmarshaler.Unmarshal(b2[:n], f.Value)
	if err != nil {
		return err
	}
	rebuild(f)
	return nil
}

func rebuild(f *Field) {
	if !f.IsCollection() {
		return
	}
	f.Items = make([]*Field, len(f.ItemsValue))
	for i := range f.ItemsValue {
		sch := f.ItemsSchema[0]
		if !f.IsColumn() {
			sch = f.ItemsSchema[i]
		}
		f.Items[i] = &Field{
			Schema: sch,
			Value:  f.ItemsValue[i],
			Index:  i,
		}
		rebuild(f.Items[i])
	}
}

// little-endian uint64
func setLen64(b2 []byte, n int, i int) {
	b2[i+0] = byte(n >> 0)
	b2[i+1] = byte(n >> 8)
	b2[i+2] = byte(n >> 16)
	b2[i+3] = byte(n >> 24)
	b2[i+4] = byte(n >> 32)
	b2[i+5] = byte(n >> 40)
	b2[i+6] = byte(n >> 48)
	b2[i+7] = byte(n >> 56)
}

// little-endian uint32
func setLen32(b2 []byte, n int, i int) {
	b2[i+0] = byte(n >> 0)
	b2[i+1] = byte(n >> 8)
	b2[i+2] = byte(n >> 16)
	b2[i+3] = byte(n >> 24)
}
