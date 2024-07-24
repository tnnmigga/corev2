package message

import "reflect"

var handlerMap = map[reflect.Type][]func(any){}

func Subscribe[T any](h func(*T)) {
	mType := reflect.TypeOf(new(T))
	handlerMap[mType] = append(handlerMap[mType], func(a any) {
		h(a.(*T))
	})
}

func Publish(msg any) {

}
