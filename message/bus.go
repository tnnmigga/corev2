package message

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mohae/deepcopy"
	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/conf"
	"github.com/tnnmigga/corev2/iface"
	"github.com/tnnmigga/corev2/infra/nmq"
	"github.com/tnnmigga/corev2/log"
	"github.com/tnnmigga/corev2/message/codec"
	"github.com/tnnmigga/corev2/utils"
	"github.com/tnnmigga/corev2/utils/conv"
)

var subMap = map[reflect.Type][]iface.IModule{}

func Handle[T any](m iface.IModule, h func(*T)) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	subscribe[T](m)
	m.Handle(mType, func(a any) {
		h(a.(*T))
	})
}

func Response[T1 any, T2 any](m iface.IModule, h func(body *T1, response func(*T2, error))) {
	codec.Register[T1]()
	mType := reflect.TypeOf(new(T1))
	subscribe[T1](m)
	m.Response(mType, func(ctx iface.IRequestCtx) {
		body := ctx.Body()
		h(body.(*T1), func(response *T2, err error) {
			ctx.Return(response, err)
		})
	})
}

func subscribe[T any](m iface.IModule) {
	mType := reflect.TypeOf(new(T))
	subMap[mType] = append(subMap[mType], m)
}

// 投递消息
func Cast(serverID uint32, msg any) error {
	if serverID == conf.ServerID {
		Delivery(msg)
	}
	b := codec.Encode(msg)
	err := nmq.Default().Publish(castSubject(serverID), b)
	if err != nil {
		log.Error(err)
	}
	return err
}

func Stream(serverID uint32, msg any) error {
	if serverID == conf.ServerID {
		Delivery(msg)
	}
	b := codec.Encode(msg)
	_, err := nmq.Default().Stream().PublishAsync(streamCastSubject(serverID), b)
	if err != nil {
		log.Error(err)
	}
	return err
}

// 投递到本地
func Delivery(msg any) {
	subs, ok := subMap[reflect.TypeOf(msg)]
	if !ok {
		log.Errorf("message cast recv not fuound %v", utils.TypeName(msg))
		return
	}
	for _, sub := range subs {
		sub.Assign(deepcopy.Copy(msg))
	}
}

// 广播到包含某个模块的所有进程
func Broadcast(name string, msg any) error {
	b := codec.Encode(msg)
	err := nmq.Default().Publish(broadcastSubject(name), b)
	if err != nil {
		log.Error(err)
	}
	return err
}

// 随机投递到一个分组下的任意进程
func Anycast(group string, msg any) error {
	b := codec.Encode(msg)
	err := nmq.Default().Publish(castAnySubject(group), b)
	if err != nil {
		log.Error(err)
	}
	return err
}

func Request[T any](serverID uint32, req any) (*T, error) {
	if serverID == conf.ServerID {
		return requestLocal[T](req)
	}
	b := codec.Encode(req)
	return request[T](requestSubject(serverID), b)
}

func RequestAsync[T any](caller iface.IModule, serverID uint32, req any, cb func(resp *T, err error)) {
	if serverID == conf.ServerID {
		conc.Go(func() {
			resp, err := requestLocal[T](req)
			caller.Assign(func() {
				cb(resp, err)
			})
		})
		return
	}
	b := codec.Encode(req)
	conc.Go(func() {
		resp, err := request[T](requestSubject(serverID), b)
		caller.Assign(func() {
			cb(resp, err)
		})
	})
}

func RequestAnyAsync[T any](caller iface.IModule, group string, req any, cb func(resp *T, err error)) {
	b := codec.Encode(req)
	conc.Go(func() {
		resp, err := request[T](requestAnySubject(group), b)
		caller.Assign(func() {
			cb(resp, err)
		})
	})
}

func RequestAny[T any](group string, req any) (*T, error) {
	b := codec.Encode(req)
	return request[T](requestAnySubject(group), b)
}

func request[T any](subject string, b []byte) (*T, error) {
	msg, err := nmq.Default().Conn.Request(subject, b, defaultTimeout)
	if err != nil {
		return nil, err
	}
	if errs := msg.Header.Values("err"); len(errs) > 0 {
		return nil, errors.New(errs[0])
	}
	result, err := codec.Decode(msg.Data)
	if err != nil {
		return nil, fmt.Errorf("request response decode error: %v", err)
	}
	data, ok := conv.Pointer[T](result)
	if !ok {
		return nil, fmt.Errorf("request response type error: %v", utils.TypeName(result))
	}
	return data, err
}

func requestLocal[T any](req any) (*T, error) {
	ctx := newRequestCtx(req)
	data, err := ctx.do()
	if err != nil {
		return nil, err
	}
	result, ok := conv.Pointer[T](data)
	if !ok {
		return nil, fmt.Errorf("RequestLocal response type error: %v", utils.TypeName(ctx.resp))
	}
	return result, err
}
