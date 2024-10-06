package message

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync/atomic"
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

func castAnySubject(group string) string {
	return fmt.Sprintf("castany.%s", group)
}

func requestSubject(serverID uint32) string {
	return fmt.Sprintf("request.%d", serverID)
}

func requestAnySubject(group string) string {
	return fmt.Sprintf("requestany.%s", group)
}

type RequestCtx struct {
	context.Context
	flag   int32
	cancel func()
	req    any
	resp   any
	err    error
}

func newRequestCtx(req any) *RequestCtx {
	ctx := &RequestCtx{
		req: deepcopy.Copy(req),
	}
	ctx.Context, ctx.cancel = context.WithTimeout(context.Background(), defaultTimeout)
	return ctx
}

func (ctx *RequestCtx) Body() any {
	return ctx.req
}

func (ctx *RequestCtx) Return(resp any, err error) {
	if !atomic.CompareAndSwapInt32(&ctx.flag, 0, 1) {
		log.Panicf("repeated return")
	}
	ctx.resp = resp
	ctx.err = err
	ctx.cancel()
}

func (ctx *RequestCtx) do() (any, error) {
	subs, ok := subMap[reflect.TypeOf(ctx.Body())]
	if !ok {
		return nil, fmt.Errorf("callee not fuound %v", utils.TypeName(ctx.Body()))
	}
	subs[0].Assign(ctx)
	<-ctx.Done()
	if ctx.Err() != nil {
		return nil, fmt.Errorf("do %v error %v", utils.TypeName(ctx.req), ctx.Err())
	}
	return ctx.resp, ctx.err
}
