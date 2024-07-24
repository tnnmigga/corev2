package codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/tnnmigga/corev2/logger"
	"github.com/tnnmigga/corev2/utils"
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
func Register[T any]() {
	msg := new(T)
	msgName := utils.TypeName(msg)
	msgID := nameToID(msgName)
	if desc, has := msgIDToDesc[msgID]; has {
		if desc.MessageName != msgName {
			logger.Panicf("msgid duplicat %v %d", msgName, msgID)
		}
	}
	msgIDToDesc[msgID] = &MessageDescriptor{
		MessageName: msgName,
		MarshalType: marshalType(msg),
		ReflectType: reflect.TypeOf(msg),
	}
}

// 编码
func Encode(msg any) []byte {
	msgID := nameToID(utils.TypeName(msg))
	bytes := Marshal(msg)
	body := make([]byte, 4, len(bytes)+4)
	binary.LittleEndian.PutUint32(body, msgID)
	body = append(body, bytes...)
	return body
}

// 解码
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
			logger.Panic(fmt.Errorf("message encode error %v", err))
		}
		return b
	}
	b, err := bson.Marshal(v)
	if err != nil {
		logger.Panic(fmt.Errorf("message encode error %v", err))
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

func nameToID(msgName string) uint32 {
	d := utils.StringToBytes(msgName)
	p := uint32(31)
	n := uint32(0)
	for _, b := range d {
		n = n*p + uint32(b)
	}
	return n
}
