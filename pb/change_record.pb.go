// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: change_record.proto

package __

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

type ChangeRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type             string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	OrderNumber      string `protobuf:"bytes,2,opt,name=order_number,json=orderNumber,proto3" json:"order_number,omitempty"`
	OrderVerb        string `protobuf:"bytes,3,opt,name=order_verb,json=orderVerb,proto3" json:"order_verb,omitempty"`
	Quantity         int64  `protobuf:"varint,4,opt,name=quantity,proto3" json:"quantity,omitempty"`
	ExecutedQuantity int64  `protobuf:"varint,5,opt,name=executed_quantity,json=executedQuantity,proto3" json:"executed_quantity,omitempty"`
	OrderBook        string `protobuf:"bytes,6,opt,name=order_book,json=orderBook,proto3" json:"order_book,omitempty"`
	Price            int64  `protobuf:"varint,7,opt,name=price,proto3" json:"price,omitempty"`
	ExecutionPrice   int64  `protobuf:"varint,8,opt,name=execution_price,json=executionPrice,proto3" json:"execution_price,omitempty"`
	StockCode        string `protobuf:"bytes,9,opt,name=stock_code,json=stockCode,proto3" json:"stock_code,omitempty"`
}

func (x *ChangeRecord) Reset() {
	*x = ChangeRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_change_record_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeRecord) ProtoMessage() {}

func (x *ChangeRecord) ProtoReflect() protoreflect.Message {
	mi := &file_change_record_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeRecord.ProtoReflect.Descriptor instead.
func (*ChangeRecord) Descriptor() ([]byte, []int) {
	return file_change_record_proto_rawDescGZIP(), []int{0}
}

func (x *ChangeRecord) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *ChangeRecord) GetOrderNumber() string {
	if x != nil {
		return x.OrderNumber
	}
	return ""
}

func (x *ChangeRecord) GetOrderVerb() string {
	if x != nil {
		return x.OrderVerb
	}
	return ""
}

func (x *ChangeRecord) GetQuantity() int64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *ChangeRecord) GetExecutedQuantity() int64 {
	if x != nil {
		return x.ExecutedQuantity
	}
	return 0
}

func (x *ChangeRecord) GetOrderBook() string {
	if x != nil {
		return x.OrderBook
	}
	return ""
}

func (x *ChangeRecord) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *ChangeRecord) GetExecutionPrice() int64 {
	if x != nil {
		return x.ExecutionPrice
	}
	return 0
}

func (x *ChangeRecord) GetStockCode() string {
	if x != nil {
		return x.StockCode
	}
	return ""
}

type ChangeRecords struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChangeRecords []*ChangeRecord `protobuf:"bytes,1,rep,name=change_records,json=changeRecords,proto3" json:"change_records,omitempty"`
}

func (x *ChangeRecords) Reset() {
	*x = ChangeRecords{}
	if protoimpl.UnsafeEnabled {
		mi := &file_change_record_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeRecords) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeRecords) ProtoMessage() {}

func (x *ChangeRecords) ProtoReflect() protoreflect.Message {
	mi := &file_change_record_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeRecords.ProtoReflect.Descriptor instead.
func (*ChangeRecords) Descriptor() ([]byte, []int) {
	return file_change_record_proto_rawDescGZIP(), []int{1}
}

func (x *ChangeRecords) GetChangeRecords() []*ChangeRecord {
	if x != nil {
		return x.ChangeRecords
	}
	return nil
}

var File_change_record_proto protoreflect.FileDescriptor

var file_change_record_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaa, 0x02, 0x0a, 0x0c, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1d, 0x0a,
	0x0a, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x56, 0x65, 0x72, 0x62, 0x12, 0x1a, 0x0a, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x2b, 0x0a, 0x11, 0x65, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x65, 0x64, 0x5f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x10, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x64, 0x51, 0x75, 0x61,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62,
	0x6f, 0x6f, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x42, 0x6f, 0x6f, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x65, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x43, 0x6f,
	0x64, 0x65, 0x22, 0x45, 0x0a, 0x0d, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x73, 0x12, 0x34, 0x0a, 0x0e, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x0d, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_change_record_proto_rawDescOnce sync.Once
	file_change_record_proto_rawDescData = file_change_record_proto_rawDesc
)

func file_change_record_proto_rawDescGZIP() []byte {
	file_change_record_proto_rawDescOnce.Do(func() {
		file_change_record_proto_rawDescData = protoimpl.X.CompressGZIP(file_change_record_proto_rawDescData)
	})
	return file_change_record_proto_rawDescData
}

var file_change_record_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_change_record_proto_goTypes = []interface{}{
	(*ChangeRecord)(nil),  // 0: ChangeRecord
	(*ChangeRecords)(nil), // 1: ChangeRecords
}
var file_change_record_proto_depIdxs = []int32{
	0, // 0: ChangeRecords.change_records:type_name -> ChangeRecord
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_change_record_proto_init() }
func file_change_record_proto_init() {
	if File_change_record_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_change_record_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeRecord); i {
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
		file_change_record_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeRecords); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_change_record_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_change_record_proto_goTypes,
		DependencyIndexes: file_change_record_proto_depIdxs,
		MessageInfos:      file_change_record_proto_msgTypes,
	}.Build()
	File_change_record_proto = out.File
	file_change_record_proto_rawDesc = nil
	file_change_record_proto_goTypes = nil
	file_change_record_proto_depIdxs = nil
}