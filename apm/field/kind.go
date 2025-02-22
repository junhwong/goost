package field

import "strings"

const (
	InvalidKind  = Type_UNKNOWN
	StringKind   = Type_STRING
	IntKind      = Type_INT
	UintKind     = Type_UINT
	FloatKind    = Type_FLOAT
	BoolKind     = Type_BOOL
	TimeKind     = Type_TIMESTAMP
	DurationKind = Type_DURATION
	BytesKind    = Type_BYTES
	IPKind       = Type_IP
	LevelKind    = Type_LOGLEVEL
	GroupKind    = Type_GROUP
	ArrayKind    = Type_ARRAY
)

func ParseType(v any) Type {
	var n int32
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
		n = int32(ParseType(Type_value[v]))
	case int:
		n = int32(v)
	case int32:
		n = v
	}
	if _, ok := Type_name[n]; ok {
		return Type(n)
	}
	return InvalidKind
}
