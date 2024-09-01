package message

import (
	fmt "fmt"
	"reflect"
	"time"

	"github.com/mohae/deepcopy"
	"github.com/tnnmigga/corev2/conf"
	"github.com/tnnmigga/corev2/utils"
)

const (
	defaultTimeout = time.Second * 10
	defaultMaxMsgs = 10000000
)

func streamName() string {
	return fmt.Sprintf("stream-cast-%d", conf.ServerID)
}

func consumerName() string {
	return fmt.Sprintf("consumer-%d", conf.ServerID)
}

func castSubject(serverID uint32) string {
	return fmt.Sprintf("cast.%d", serverID)
}

func streamCastSubject(serverID uint32) string {
	return fmt.Sprintf("stream.cast.%d", serverID)
}

func broadcastSubject(group string) string {
	return fmt.Sprintf("broadcast.%s", group)
}

func anycastSubject(group string) string {
	return fmt.Sprintf("randomcast.%s", group)
}

func rpcSubject(serverID uint32) string {
	return fmt.Sprintf("rpc.%d", serverID)
}

func anyRPCSubject(group string) string {
	return fmt.Sprintf("randomrpc.%s", group)
}

type RPCContext struct {
	req  any
	resp any
	err  error
	cb   func(resp any, err error)
	sign chan any
}

func newRPCContext(req any) *RPCContext {
	ctx := &RPCContext{
		req:  deepcopy.Copy(req),
		sign: make(chan any, 1),
	}
	ctx.cb = func(_resp any, _err error) {
		ctx.err = _err
		if _err != nil {
			ctx.sign <- struct{}{}
			return
		}
		ctx.resp = _resp
		ctx.sign <- struct{}{}
	}
	return ctx
}

func (ctx *RPCContext) RPCBody() any {
	return ctx.req
}

func (ctx *RPCContext) Return(resp any, err error) {
	ctx.cb(resp, err)
}

func (ctx *RPCContext) exec() (any, error) {
	subs, ok := subMap[reflect.TypeOf(ctx.req)]
	if !ok {
		return nil, fmt.Errorf("localCall callee not fuound %v", utils.TypeName(ctx.req))
	}
	subs[0].Assign(ctx)
	timeout := time.NewTimer(defaultTimeout)
	defer timeout.Stop()
	select {
	case <-ctx.sign:
		return ctx.resp, ctx.err
	case <-timeout.C:
		return nil, fmt.Errorf("localCall %v timeout", utils.TypeName(ctx.req))
	}
}
