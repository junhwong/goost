package field

import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"

	"math/bits"

	"github.com/junhwong/goost/buffering"
)

func Marshal(f *Field, buf *buffering.Buffer, keepName bool) (err error) {
	if f == nil || f.kind == InvalidKind {
		return nil
	}

	oldoff := buf.Len()
	defer func() {
		if err == nil {
			return
		}
		buf.Truncate(oldoff)
	}()

	// ver type flags name subtypes payload
	write := func(b []byte) {
		if err != nil {
			return
		}
		_, err = buf.Write(b)
	}
	writeStr := func(s string) {
		b := *(*[]byte)(unsafe.Pointer(&s))
		write(AppendVarint(nil, uint64(len(b))))
		write(b)
	}
	writeInt := func(i int64) {
		if err != nil {
			return
		}
		write(AppendVarint(nil, EncodeZigZag(i)))
	}

	var writeType func(f *Field, kn bool)
	writeType = func(f *Field, kn bool) {
		if err != nil {
			return
		}
		err = buf.WriteByte(byte(f.kind))
		if err != nil {
			return
		}
		if e := binary.Write(buf, binary.BigEndian, f.typFlag); e != nil {
			err = e
			return
		}

		if kn {
			writeStr(f.GetName())
		} else {
			writeStr("")
		}

		switch f.kind {
		case ListKind, DictKind:
			write(AppendVarint(nil, uint64(len(f.Items))))
			for _, v := range f.Items {
				writeType(v, f.kind == DictKind)
			}
		case StringKind, IntKind, UintKind, FloatKind, BoolKind, TimeKind, DurationKind, BytesKind, IPKind, LevelKind:
		default:
			err = fmt.Errorf("unsupported type %v", f.kind)
		}
	}
	var writeVal func(f *Field)
	writeVal = func(f *Field) {
		if err != nil {
			return
		}
		if e := binary.Write(buf, binary.BigEndian, f.valFlag); e != nil {
			err = e
			return
		}
		if f.IsNull() {
			return
		}
		switch f.kind {
		case StringKind:
			writeStr(f.GetString())
		case IntKind:
			writeInt(f.GetInt())
		case UintKind:
			write(AppendVarint(nil, f.GetUint()))
		case FloatKind:
			// 补充浮点处理（按 IEEE754 双精度写入）
			var buf [8]byte
			binary.BigEndian.PutUint64(buf[:], math.Float64bits(f.GetFloat()))
			write(buf[:])
		case BoolKind:
			if f.GetBool() {
				err = buf.WriteByte(1)
			} else {
				err = buf.WriteByte(0)
			}
		case TimeKind:
			t := f.GetTime()
			writeInt(t.UnixNano())
		case DurationKind:
			writeInt(f.GetDuration().Nanoseconds())
		case BytesKind:
			b := f.GetBytes()
			write(AppendVarint(nil, uint64(len(b))))
			write(b)
		case IPKind:
			ip := f.GetIP()
			write(AppendVarint(nil, uint64(len(ip))))
			write(ip)
		case LevelKind:
			write(AppendVarint(nil, uint64(f.GetLevel())))
		case ListKind, DictKind:
			write(AppendVarint(nil, uint64(len(f.Items))))
			for _, v := range f.Items {
				writeVal(v)
			}
		}
	}

	err = buf.WriteByte(0) // ver
	writeType(f, keepName)
	writeVal(f)
	return
}
func Unmarshal(buf *buffering.Buffer, f *Field) (err error) {
	// 读取版本号
	ver, e := buf.ReadByte()
	if e != nil {
		return e
	}
	if ver != 0 {
		return fmt.Errorf("unsupported version: %d", ver)
	}

	readLen := func() int {
		if err != nil {
			return -1
		}
		v, n := ConsumeVarint(buf.Bytes())
		if n < 0 {
			return -1
		}
		if n > 5 {
			panic("unsupported length")
		}
		buf.Next(n)
		return int(v)
	}

	// 递归读取类型信息
	var readType func(parent *Field)
	readType = func(parent *Field) {
		if err != nil {
			return
		}
		// 读取类型和标志位
		typ, e := buf.ReadByte()
		if e != nil {
			err = e
			return
		}

		parent.kind = Type(typ)
		if e := binary.Read(buf, binary.BigEndian, &parent.typFlag); e != nil {
			err = e
			return
		}

		// 读取名称
		n := readLen()
		if n < 0 {
			return
		}
		if n != 0 {
			b := make([]byte, n)
			if _, e := buf.Read(b); e != nil {
				err = e
				return
			}
			parent.name = string(b)
		}

		// 处理子类型
		switch parent.kind {
		case ListKind, DictKind:
			n := readLen()
			if n < 0 {
				return
			}

			parent.Items = make([]*Field, n)
			for i := 0; i < n; i++ {
				child := fieldPool.Get().(*Field)
				parent.Items[i] = child
				readType(child)
			}
		case StringKind, IntKind, UintKind, FloatKind, BoolKind, TimeKind, DurationKind, BytesKind, IPKind, LevelKind:
			// 基础类型无需子项
		default:
			err = fmt.Errorf("unsupported type %v", parent.kind)
			return
		}
	}
	readVarInt := func() uint64 {
		if err != nil {
			return 0
		}
		v, n := ConsumeVarint(buf.Bytes())
		if n < 0 {
			err = fmt.Errorf("read varint err: %v", e)
		}
		buf.Next(n)
		return v
	}
	readInt := func() int64 {
		v := readVarInt()
		return DecodeZigZag(v)
	}

	// 递归读取值
	var readVal func(f *Field)
	readVal = func(f *Field) {
		if err != nil {
			return
		}
		// 读取标志位
		if e := binary.Read(buf, binary.BigEndian, &f.valFlag); e != nil {
			err = e
			return
		}
		// 空值直接返回
		if f.IsNull() {
			return
		}

		switch f.kind {
		case StringKind:
			n := readLen()
			if n < 0 {
				return
			}
			if n != 0 {
				b := make([]byte, n)
				if _, e := buf.Read(b); e != nil {
					err = e
					return
				}
				f.strVal = string(b)
			}
		case IntKind, DurationKind, LevelKind, TimeKind:
			f.intVal = readInt()
		case UintKind:
			f.uintVal = readVarInt()
		case FloatKind:
			var b [8]byte
			if _, e := buf.Read(b[:]); e != nil {
				err = e
				return
			}
			f.floatVal = math.Float64frombits(binary.BigEndian.Uint64(b[:]))

		case BoolKind:
			b, e := buf.ReadByte()
			if e != nil {
				err = e
				return
			}
			f.intVal = int64(b)
		case BytesKind, IPKind:
			n := readLen()
			if n < 0 {
				return
			}
			f.bytesVal = make([]byte, n)
			if _, e := buf.Read(f.bytesVal); e != nil {
				err = e
				return
			}
		case ListKind, DictKind:
			count := readLen()
			if count < 0 {
				return
			}
			if len(f.Items) < count {
				old := f.Items
				f.Items = make([]*Field, count)
				for i := 0; i < len(old); i++ {
					f.Items[i] = old[i]
				}
			}
			for i := 0; i < int(count); i++ {
				child := f.Items[i]
				if child == nil {
					if i == 0 {
						panic("丢失metadata")
					}
					prev := f.Items[i-1]
					if prev == nil {
						panic("丢失metadata")
					}
					child = fieldPool.Get().(*Field)
					child.name = prev.name
					CloneInto(prev, child)
					f.Items[i] = child
				}

				readVal(child)
			}
		default:
			err = fmt.Errorf("unsupported value type %v", f.kind)
			return
		}
	}

	readType(f)
	readVal(f)
	return
}

// func rebuild(f *Field) {
// 	if !f.IsCollection() {
// 		return
// 	}
// 	f.Items = make([]*Field, len(f.val.ItemsValue))
// 	for i := range f.val.ItemsValue {
// 		sch := f.sch.ItemsSchema[0]
// 		if !f.IsColumn() {
// 			sch = f.sch.ItemsSchema[i]
// 		}
// 		f.Items[i] = &Field{
// 			sch:   sch,
// 			val:   f.val.ItemsValue[i],
// 			Index: i,
// 		}
// 		rebuild(f.Items[i])
// 	}
// }

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

// AppendVarint appends v to b as a varint-encoded uint64.
func AppendVarint(b []byte, v uint64) []byte {
	switch {
	case v < 1<<7:
		b = append(b, byte(v))
	case v < 1<<14:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte(v>>7))
	case v < 1<<21:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte(v>>14))
	case v < 1<<28:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte(v>>21))
	case v < 1<<35:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte(v>>28))
	case v < 1<<42:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte(v>>35))
	case v < 1<<49:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte(v>>42))
	case v < 1<<56:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte(v>>49))
	case v < 1<<63:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte(v>>56))
	default:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte((v>>56)&0x7f|0x80),
			1)
	}
	return b
}

const (
	_ = -iota
	errCodeTruncated
	errCodeFieldNumber
	errCodeOverflow
	errCodeReserved
	errCodeEndGroup
	errCodeRecursionDepth
)

// ConsumeVarint parses b as a varint-encoded uint64, reporting its length.
// This returns a negative length upon an error (see [ParseError]).
func ConsumeVarint(b []byte) (v uint64, n int) {
	var y uint64
	if len(b) <= 0 {
		return 0, errCodeTruncated
	}
	v = uint64(b[0])
	if v < 0x80 {
		return v, 1
	}
	v -= 0x80

	if len(b) <= 1 {
		return 0, errCodeTruncated
	}
	y = uint64(b[1])
	v += y << 7
	if y < 0x80 {
		return v, 2
	}
	v -= 0x80 << 7

	if len(b) <= 2 {
		return 0, errCodeTruncated
	}
	y = uint64(b[2])
	v += y << 14
	if y < 0x80 {
		return v, 3
	}
	v -= 0x80 << 14

	if len(b) <= 3 {
		return 0, errCodeTruncated
	}
	y = uint64(b[3])
	v += y << 21
	if y < 0x80 {
		return v, 4
	}
	v -= 0x80 << 21

	if len(b) <= 4 {
		return 0, errCodeTruncated
	}
	y = uint64(b[4])
	v += y << 28
	if y < 0x80 {
		return v, 5
	}
	v -= 0x80 << 28

	if len(b) <= 5 {
		return 0, errCodeTruncated
	}
	y = uint64(b[5])
	v += y << 35
	if y < 0x80 {
		return v, 6
	}
	v -= 0x80 << 35

	if len(b) <= 6 {
		return 0, errCodeTruncated
	}
	y = uint64(b[6])
	v += y << 42
	if y < 0x80 {
		return v, 7
	}
	v -= 0x80 << 42

	if len(b) <= 7 {
		return 0, errCodeTruncated
	}
	y = uint64(b[7])
	v += y << 49
	if y < 0x80 {
		return v, 8
	}
	v -= 0x80 << 49

	if len(b) <= 8 {
		return 0, errCodeTruncated
	}
	y = uint64(b[8])
	v += y << 56
	if y < 0x80 {
		return v, 9
	}
	v -= 0x80 << 56

	if len(b) <= 9 {
		return 0, errCodeTruncated
	}
	y = uint64(b[9])
	v += y << 63
	if y < 2 {
		return v, 10
	}
	return 0, errCodeOverflow
}

// SizeVarint returns the encoded size of a varint.
// The size is guaranteed to be within 1 and 10, inclusive.
func SizeVarint(v uint64) int {
	// This computes 1 + (bits.Len64(v)-1)/7.
	// 9/64 is a good enough approximation of 1/7
	return int(9*uint32(bits.Len64(v))+64) / 64
}

// DecodeZigZag decodes a zig-zag-encoded uint64 as an int64.
//
//	Input:  {…,  5,  3,  1,  0,  2,  4,  6, …}
//	Output: {…, -3, -2, -1,  0, +1, +2, +3, …}
func DecodeZigZag(x uint64) int64 {
	return int64(x>>1) ^ int64(x)<<63>>63
}

// EncodeZigZag encodes an int64 as a zig-zag-encoded uint64.
//
//	Input:  {…, -3, -2, -1,  0, +1, +2, +3, …}
//	Output: {…,  5,  3,  1,  0,  2,  4,  6, …}
func EncodeZigZag(x int64) uint64 {
	return uint64(x<<1) ^ uint64(x>>63)
}
