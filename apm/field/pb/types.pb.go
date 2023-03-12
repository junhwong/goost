// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        v4.22.1
// source: types.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 值类型
type Field_ValueType int32

const (
	Field_UNKNOWN   Field_ValueType = 0 // 未知
	Field_STRING    Field_ValueType = 1 // 字符串
	Field_FLOAT     Field_ValueType = 2 // 浮点数
	Field_BOOL      Field_ValueType = 3 // 布尔值
	Field_INT       Field_ValueType = 4 // 整数 int64
	Field_UINT      Field_ValueType = 5 // 整数 uint64
	Field_TIMESTAMP Field_ValueType = 6 // 时间戳 纳秒
	Field_DURATION  Field_ValueType = 7 // 持续时间 纳秒
)

// Enum value maps for Field_ValueType.
var (
	Field_ValueType_name = map[int32]string{
		0: "UNKNOWN",
		1: "STRING",
		2: "FLOAT",
		3: "BOOL",
		4: "INT",
		5: "UINT",
		6: "TIMESTAMP",
		7: "DURATION",
	}
	Field_ValueType_value = map[string]int32{
		"UNKNOWN":   0,
		"STRING":    1,
		"FLOAT":     2,
		"BOOL":      3,
		"INT":       4,
		"UINT":      5,
		"TIMESTAMP": 6,
		"DURATION":  7,
	}
)

func (x Field_ValueType) Enum() *Field_ValueType {
	p := new(Field_ValueType)
	*p = x
	return p
}

func (x Field_ValueType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Field_ValueType) Descriptor() protoreflect.EnumDescriptor {
	return file_types_proto_enumTypes[0].Descriptor()
}

func (Field_ValueType) Type() protoreflect.EnumType {
	return &file_types_proto_enumTypes[0]
}

func (x Field_ValueType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Field_ValueType.Descriptor instead.
func (Field_ValueType) EnumDescriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{0, 0}
}

// 标签
type Field struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key         string          `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`                             // 名称
	Kind        Field_ValueType `protobuf:"varint,2,opt,name=kind,proto3,enum=apm.Field_ValueType" json:"kind,omitempty"` // 类型
	StringValue *string         `protobuf:"bytes,4,opt,name=stringValue,proto3,oneof" json:"stringValue,omitempty"`
	BoolValue   *bool           `protobuf:"varint,5,opt,name=boolValue,proto3,oneof" json:"boolValue,omitempty"`
	IntValue    *int64          `protobuf:"varint,6,opt,name=intValue,proto3,oneof" json:"intValue,omitempty"`
	UintValue   *uint64         `protobuf:"varint,7,opt,name=uintValue,proto3,oneof" json:"uintValue,omitempty"`
	FloatValue  *float64        `protobuf:"fixed64,8,opt,name=floatValue,proto3,oneof" json:"floatValue,omitempty"`
	IsTag       *bool           `protobuf:"varint,11,opt,name=isTag,proto3,oneof" json:"isTag,omitempty"` // 是否是索引
}

func (x *Field) Reset() {
	*x = Field{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{0}
}

func (x *Field) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Field) GetKind() Field_ValueType {
	if x != nil {
		return x.Kind
	}
	return Field_UNKNOWN
}

func (x *Field) GetStringValue() string {
	if x != nil && x.StringValue != nil {
		return *x.StringValue
	}
	return ""
}

func (x *Field) GetBoolValue() bool {
	if x != nil && x.BoolValue != nil {
		return *x.BoolValue
	}
	return false
}

func (x *Field) GetIntValue() int64 {
	if x != nil && x.IntValue != nil {
		return *x.IntValue
	}
	return 0
}

func (x *Field) GetUintValue() uint64 {
	if x != nil && x.UintValue != nil {
		return *x.UintValue
	}
	return 0
}

func (x *Field) GetFloatValue() float64 {
	if x != nil && x.FloatValue != nil {
		return *x.FloatValue
	}
	return 0
}

func (x *Field) GetIsTag() bool {
	if x != nil && x.IsTag != nil {
		return *x.IsTag
	}
	return false
}

var File_types_proto protoreflect.FileDescriptor

var file_types_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61,
	0x70, 0x6d, 0x22, 0xce, 0x03, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x28,
	0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x61,
	0x70, 0x6d, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x25, 0x0a, 0x0b, 0x73, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x21, 0x0a, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x01, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x88,
	0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x75, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x48, 0x03, 0x52, 0x09, 0x75, 0x69, 0x6e, 0x74, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x0a, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x48, 0x04, 0x52, 0x0a, 0x66, 0x6c,
	0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x69,
	0x73, 0x54, 0x61, 0x67, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x48, 0x05, 0x52, 0x05, 0x69, 0x73,
	0x54, 0x61, 0x67, 0x88, 0x01, 0x01, 0x22, 0x69, 0x0a, 0x09, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05,
	0x46, 0x4c, 0x4f, 0x41, 0x54, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x42, 0x4f, 0x4f, 0x4c, 0x10,
	0x03, 0x12, 0x07, 0x0a, 0x03, 0x49, 0x4e, 0x54, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x55, 0x49,
	0x4e, 0x54, 0x10, 0x05, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x49, 0x4d, 0x45, 0x53, 0x54, 0x41, 0x4d,
	0x50, 0x10, 0x06, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x55, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10,
	0x07, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42,
	0x0b, 0x0a, 0x09, 0x5f, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x0c, 0x0a, 0x0a,
	0x5f, 0x75, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x66,
	0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x69, 0x73,
	0x54, 0x61, 0x67, 0x42, 0x0e, 0x5a, 0x0c, 0x61, 0x70, 0x6d, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_types_proto_rawDescOnce sync.Once
	file_types_proto_rawDescData = file_types_proto_rawDesc
)

func file_types_proto_rawDescGZIP() []byte {
	file_types_proto_rawDescOnce.Do(func() {
		file_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_types_proto_rawDescData)
	})
	return file_types_proto_rawDescData
}

var file_types_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_types_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_types_proto_goTypes = []interface{}{
	(Field_ValueType)(0), // 0: apm.Field.ValueType
	(*Field)(nil),        // 1: apm.Field
}
var file_types_proto_depIdxs = []int32{
	0, // 0: apm.Field.kind:type_name -> apm.Field.ValueType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_types_proto_init() }
func file_types_proto_init() {
	if File_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Field); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_types_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_types_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_types_proto_goTypes,
		DependencyIndexes: file_types_proto_depIdxs,
		EnumInfos:         file_types_proto_enumTypes,
		MessageInfos:      file_types_proto_msgTypes,
	}.Build()
	File_types_proto = out.File
	file_types_proto_rawDesc = nil
	file_types_proto_goTypes = nil
	file_types_proto_depIdxs = nil
}
