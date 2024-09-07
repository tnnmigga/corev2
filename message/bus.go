package message

import (
	"errors"
	fmt "fmt"
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

func Subscribe[T any](m iface.IModule) {
	mType := reflect.TypeOf(new(T))
	subMap[mType] = append(subMap[mType], m)
}

func Handle[T any](m iface.IModule, h func(*T)) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	Subscribe[T](m)
	m.Handle(mType, func(a any) {
		h(a.(*T))
	})
}

func RegisterRPC[T any](m iface.IModule, rpc func(req *T, resp func(any, error))) {
	codec.Register[T]()
	mType := reflect.TypeOf(new(T))
	Subscribe[T](m)
	m.RegisterRPC(mType, func(req iface.IRPCCtx) {
		body := req.RPCBody()
		rpc(body.(*T), req.Return)
	})
}

// 跨进程投递消息
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

// 广播到一个分组下的所有进程
func Broadcast(group string, msg any) error {
	b := codec.Encode(msg)
	err := nmq.Default().Publish(broadcastSubject(group), b)
	if err != nil {
		log.Error(err)
	}
	return err
}

// 随机投递到一个分组下的某个进程
func Anycast(group string, msg any) error {
	b := codec.Encode(msg)
	err := nmq.Default().Publish(anycastSubject(group), b)
	if err != nil {
		log.Error(err)
	}
	return err
}

// RPCAsync 跨协程/进程调用
// caller: 为调用者模块 也是回调函数的执行者
// serverID: 目标参数 可以通过msgbus.ServerID()指定某个特定的进程或通过msgbus.ServerType()在某类进程中随机一个
// 调用本地使用msgbus.Local()或msgbus.ServerID(conf.ServerID)
// req: 请求参数
// cb: 回调函数 由调用方模块线程执行
func RPCAsync[T any](caller iface.IModule, serverID uint32, req any, cb func(resp *T, err error)) {
	b := codec.Encode(req)
	conc.Go(func() {
		msg, err := nmq.Default().Conn.Request(rpcSubject(serverID), b, defaultTimeout)
		if err != nil {
			caller.Assign(func() {
				cb(nil, err)
			})
			return
		}
		if errs := msg.Header.Values("err"); len(errs) > 0 {
			caller.Assign(func() {
				cb(nil, errors.New(errs[0]))
			})
			return
		}
		result, err := codec.Decode(msg.Data)
		if err != nil {
			caller.Assign(func() {
				cb(nil, fmt.Errorf("RPCAsync response decode error: %v", err))
			})
			return
		}
		data, ok := conv.Pointer[T](result)
		if !ok {
			caller.Assign(func() {
				cb(nil, fmt.Errorf("RPCAsync response type error: %v", utils.TypeName(result)))
			})
			return
		}
		caller.Assign(func() {
			cb(data, nil)
		})
	})
}

func RPC[T any](serverID uint32, req any) (*T, error) {
	b := codec.Encode(req)
	msg, err := nmq.Default().Conn.Request(rpcSubject(serverID), b, defaultTimeout)
	if err != nil {
		return nil, err
	}
	if errs := msg.Header.Values("err"); len(errs) > 0 {
		return nil, errors.New(errs[0])
	}
	result, err := codec.Decode(msg.Data)
	if err != nil {
		return nil, fmt.Errorf("RPC response decode error: %v", err)
	}
	data, ok := conv.Pointer[T](result)
	if !ok {
		return nil, fmt.Errorf("RPC response type error: %v", utils.TypeName(result))
	}
	return data, err
}

func AnyRPCAsync[T any](caller iface.IModule, group string, req any, cb func(resp *T, err error)) {
	b := codec.Encode(req)
	conc.Go(func() {
		msg, err := nmq.Default().Conn.Request(anyRPCSubject(group), b, defaultTimeout)
		if err != nil {
			caller.Assign(func() {
				cb(nil, err)
			})
			return
		}
		if errs := msg.Header.Values("err"); len(errs) > 0 {
			caller.Assign(func() {
				cb(nil, errors.New(errs[0]))
			})
			return
		}
		result, err := codec.Decode(msg.Data)
		if err != nil {
			caller.Assign(func() {
				cb(nil, fmt.Errorf("AnyRPCAsync response decode error: %v", err))
			})
			return
		}
		data, ok := conv.Pointer[T](result)
		if !ok {
			caller.Assign(func() {
				cb(nil, fmt.Errorf("AnyRPCAsync response type error: %v", utils.TypeName(result)))
			})
			return
		}
		caller.Assign(func() {
			cb(data, nil)
		})
	})
}

func AnyRPC[T any](group string, req any) (*T, error) {
	b := codec.Encode(req)
	msg, err := nmq.Default().Conn.Request(anyRPCSubject(group), b, defaultTimeout)
	if err != nil {
		return nil, err
	}
	if errs := msg.Header.Values("err"); len(errs) > 0 {
		return nil, errors.New(errs[0])
	}
	result, err := codec.Decode(msg.Data)
	if err != nil {
		return nil, fmt.Errorf("RPC response decode error: %v", err)
	}
	data, ok := conv.Pointer[T](result)
	if !ok {
		return nil, fmt.Errorf("RPC response type error: %v", utils.TypeName(result))
	}
	return data, err
}

func LPC[T any](req any) (*T, error) {
	ctx := newRPCContext(req)
	data, err := ctx.exec()
	if err != nil {
		return nil, err
	}
	result, ok := conv.Pointer[T](data)
	if !ok {
		return nil, fmt.Errorf("LPC response type error: %v", utils.TypeName(ctx.resp))
	}
	return result, err
}

func LPCAsync[T any](caller iface.IModule, req any, cb func(resp *T, err error)) {
	ctx := newRPCContext(req)
	conc.Go(func() {
		data, err := ctx.exec()
		if err != nil {
			caller.Assign(func() {
				cb(nil, err)
			})
			return
		}
		result, ok := conv.Pointer[T](data)
		if !ok {
			caller.Assign(func() {
				cb(nil, fmt.Errorf("LPCAsync response type error: %v", utils.TypeName(ctx.resp)))
			})
			return
		}
		caller.Assign(func() {
			cb(result, nil)
		})
	})
}
