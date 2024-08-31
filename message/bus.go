package message

import (
	"reflect"

	"github.com/tnnmigga/corev2/iface"
)

var handlerMap = map[reflect.Type][]func(any){}

func Subscribe[T any](m iface.IModule) {
	mType := reflect.TypeOf(new(T))
	handlerMap[mType] = append(handlerMap[mType], func(a any) {
		m.Assign(a)
	})
}

func Publish(msg any) {
	
}
