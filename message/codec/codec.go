package codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/tnnmigga/corev2/utils/stack"
	"github.com/tnnmigga/corev2/zlog"

	"github.com/gogo/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	msgIDToDesc map[uint32]*MessageDescriptor
)

func init() {
	msgIDToDesc = map[uint32]*MessageDescriptor{}
}

const (
	marshalTypeGogoproto = iota
	marshalTypeBSON
)

type MessageDescriptor struct {
	MessageName string
	MarshalType int
	ReflectType reflect.Type
}

func (d *MessageDescriptor) New() any {
	return reflect.New(d.ReflectType).Interface()
}

// 将消息注册到解码器
// 可以通过指定泛型类型和参数传递类型
func Register[T any](t ...T) {
	if len(t) == 0 {
		var tmp T
		t = append(t, tmp)
	}
	v := t[0]
	name := stack.TypeName(v)
	id := stack.TypeID(v)
	if desc, has := msgIDToDesc[id]; has {
		if desc.MessageName != name {
			zlog.Panicf("msgid duplicat %v %d", name, id)
		}
	}
	mType := reflect.TypeOf(v)
	if mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}
	msgIDToDesc[id] = &MessageDescriptor{
		MessageName: name,
		MarshalType: marshalType(v),
		ReflectType: mType,
	}
}

// 编码
// 额外拼接四字节类型id
func Encode(v any) []byte {
	msgID := stack.TypeID(v)
	bytes := Marshal(v)
	body := make([]byte, 4, len(bytes)+4)
	binary.LittleEndian.PutUint32(body, msgID)
	body = append(body, bytes...)
	return body
}

// 解码
// 使用前需要提前注册
// 需要头部四字节为类型id
func Decode(b []byte) (msg any, err error) {
	if len(b) < 4 {
		return nil, fmt.Errorf("message decode len error %d", len(b))
	}
	msgID := binary.LittleEndian.Uint32(b)
	desc, ok := msgIDToDesc[msgID]
	if !ok {
		return nil, fmt.Errorf("message decode msgid not found %d", msgID)
	}
	msg = desc.New()
	err = Unmarshal(b[4:], msg)
	return msg, err
}

// 序列化
func Marshal(v any) []byte {
	if v0, ok := v.(proto.Message); ok {
		b, err := proto.Marshal(v0)
		if err != nil {
			zlog.Panic(fmt.Errorf("message encode error %v", err))
		}
		return b
	}
	b, err := bson.Marshal(v)
	if err != nil {
		zlog.Panic(fmt.Errorf("message encode error %v", err))
	}
	return b
}

// 反序列化
// 使用前需要提前注册
func Unmarshal(b []byte, addr any) error {
	switch marshalType(addr) {
	case marshalTypeGogoproto:
		return proto.Unmarshal(b, addr.(proto.Message))
	case marshalTypeBSON:
		return bson.Unmarshal(b, addr)
	default:
		return errors.New("invalid marshal type")
	}
}

func marshalType(v any) int {
	if _, ok := v.(proto.Message); ok {
		return marshalTypeGogoproto
	}
	return marshalTypeBSON
}
