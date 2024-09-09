package codec

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/tnnmigga/corev2/log"
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
	MarshalByGoGoProto = iota
	MarshalByBSON
	MarshalByJSON
)

type MessageDescriptor struct {
	MessageName string
	ReflectType reflect.Type
}

func (d *MessageDescriptor) New() any {
	return reflect.New(d.ReflectType).Interface()
}

// 将消息注册到解码器
func Register[T any]() {
	var msg T
	msgName := utils.TypeName(msg)
	msgID := nameToID(msgName)
	if desc, has := msgIDToDesc[msgID]; has {
		if desc.MessageName != msgName {
			log.Panicf("msgid duplicat %v %d", msgName, msgID)
		}
	}
	mType := reflect.TypeOf(msg)
	if mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}
	msgIDToDesc[msgID] = &MessageDescriptor{
		MessageName: msgName,
		ReflectType: mType,
	}
}

// 编码
func Encode(msg any) []byte {
	msgID := nameToID(utils.TypeName(msg))
	mType := marshalType(msg)
	bytes := Marshal(mType, msg)
	body := make([]byte, 4, len(bytes)+5)
	binary.LittleEndian.PutUint32(body, msgID)
	body = append(append(body, mType), bytes...)
	return body
}

// 解码
func Decode(b []byte) (msg any, err error) {
	if len(b) < 5 {
		return nil, fmt.Errorf("message decode len error %d", len(b))
	}
	msgID := binary.LittleEndian.Uint32(b)
	desc, ok := msgIDToDesc[msgID]
	if !ok {
		return nil, fmt.Errorf("message decode msgid not found %d", msgID)
	}
	msg = desc.New()
	mType := b[4]
	err = Unmarshal(mType, b[5:], msg)
	return msg, err
}

// 序列化
func Marshal(mType byte, v any) []byte {
	switch mType {
	case MarshalByGoGoProto:
		b, err := proto.Marshal(v.(proto.Message))
		if err != nil {
			log.Panic(fmt.Errorf("message encode error %v", err))
		}
		return b
	case MarshalByBSON:
		b, err := bson.Marshal(v)
		if err != nil {
			log.Panic(fmt.Errorf("message encode error %v", err))
		}
		return b
	case MarshalByJSON:
		b, err := json.Marshal(v)
		if err != nil {
			log.Panic(fmt.Errorf("message encode error %v", err))
		}
		return b
	}
	log.Panic(fmt.Errorf("error marshal type %d", mType))
	return nil
}

// 反序列化
// 使用前需要提前注册
func Unmarshal(mType byte, b []byte, addr any) error {
	switch mType {
	case MarshalByGoGoProto:
		return proto.Unmarshal(b, addr.(proto.Message))
	case MarshalByBSON:
		return bson.Unmarshal(b, addr)
	default:
		return errors.New("invalid marshal type")
	}
}

func marshalType(v any) byte {
	if _, ok := v.(proto.Message); ok {
		return MarshalByGoGoProto
	}
	return MarshalByBSON
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
