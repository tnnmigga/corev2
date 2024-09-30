package message

import (
	"context"
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
	return fmt.Sprintf("anycast.%s", group)
}

func requestSubject(serverID uint32) string {
	return fmt.Sprintf("request.%d", serverID)
}

func requestAnySubject(group string) string {
	return fmt.Sprintf("requestany.%s", group)
}

type ReqCtx struct {
	context.Context
	cancel func()
	req    any
	resp   any
	err    error
}

func newReqCtx(req any) *ReqCtx {
	ctx := &ReqCtx{
		req: deepcopy.Copy(req),
	}
	ctx.Context, ctx.cancel = context.WithTimeout(context.Background(), defaultTimeout)
	return ctx
}

func (ctx *ReqCtx) ReqBody() any {
	return ctx.req
}

func (ctx *ReqCtx) Return(resp any, err error) {
	ctx.resp = resp
	ctx.err = err
	ctx.cancel()
}

func (ctx *ReqCtx) do() (any, error) {
	subs, ok := subMap[reflect.TypeOf(ctx.ReqBody())]
	if !ok {
		return nil, fmt.Errorf("callee not fuound %v", utils.TypeName(ctx.ReqBody()))
	}
	subs[0].Assign(ctx)
	<-ctx.Done()
	if ctx.Err() != nil {
		return nil, fmt.Errorf("do %v error %v", utils.TypeName(ctx.req), ctx.Err())
	}
	return ctx.resp, ctx.err
}
