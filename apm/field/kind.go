package field

import (
	"strconv"
	"strings"
)

type Type byte

const (
	InvalidKind Type = iota
	DictKind
	ListKind
	StringKind
	IntKind
	UintKind
	FloatKind
	BoolKind
	TimeKind
	DurationKind
	BytesKind
	IPKind
	LevelKind
)

func ParseType(v any) Type {
	var n Type
	switch v := v.(type) {
	case string:
		v = strings.ToUpper(v)
		switch v {
		case "TIME", "DATETIME", "DATE":
			v = "TIMESTAMP"
		case "LONG", "SHORT", "INTEGER":
			v = "INT"
		case "NUMBER", "DOUBLE":
			v = "FLOAT"
		case "STR":
			v = "STRING"
		}
		n = Type_value[v]
	case int:
		n = Type(v)
	case int32:
		n = Type(v)
	default:
		n = 0
	}
	if _, ok := Type_name[n]; ok {
		return Type(n)
	}
	return InvalidKind
}

// Enum value maps for Type.
var (
	Type_name = map[Type]string{
		0:            "UNKNOWN",
		DictKind:     "DICT",
		ListKind:     "LIST",
		StringKind:   "STRING",
		BoolKind:     "BOOL",
		IntKind:      "INT",
		UintKind:     "UINT",
		FloatKind:    "FLOAT",
		TimeKind:     "TIMESTAMP",
		DurationKind: "DURATION",
		BytesKind:    "BYTES",
		LevelKind:    "LOGLEVEL",
		IPKind:       "IP",
	}
	Type_value = map[string]Type{
		"UNKNOWN":   0,
		"MAP":       DictKind,
		"GROUP":     DictKind,
		"ARRAY":     ListKind,
		"DICT":      DictKind,
		"LIST":      ListKind,
		"STRING":    StringKind,
		"BOOL":      BoolKind,
		"INT":       IntKind,
		"UINT":      UintKind,
		"FLOAT":     FloatKind,
		"TIMESTAMP": TimeKind,
		"DURATION":  DurationKind,
		"BYTES":     BytesKind,
		"LOGLEVEL":  LevelKind,
		"IP":        IPKind,
	}
)

func (x Type) String() string {
	s, ok := Type_name[x]
	if !ok {
		return "Type(" + strconv.FormatInt(int64(x), 10) + ")"
	}
	return s
}
