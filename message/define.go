package message

import (
	fmt "fmt"
	"time"

	"github.com/mohae/deepcopy"
	"github.com/tnnmigga/corev2/utils"
)

func castSubject(serverID uint32) string {
	return fmt.Sprintf("cast.%d", serverID)
}

func streamCastSubject(serverID uint32) string {
	return fmt.Sprintf("stream.cast.%d", serverID)
}

func broadcastSubject(group string) string {
	return fmt.Sprintf("broadcast.%s", group)
}

func randomCastSubject(group string) string {
	return fmt.Sprintf("randomcast.%s", group)
}

func rpcSubject(serverID uint32) string {
	return fmt.Sprintf("rpc.%d", serverID)
}

func randomRpcSubject(serverType string) string {
	return fmt.Sprintf("randomrpc.%s", serverType)
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

func (ctx *RPCContext) wait() error {
	timeout := time.NewTimer(time.Second * 10)
	defer timeout.Stop()
	select {
	case <-ctx.sign:
		return nil
	case <-timeout.C:
		return fmt.Errorf("request %v timeout", utils.TypeName(ctx.req))
	}
}
