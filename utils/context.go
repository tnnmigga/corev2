package utils

import "context"

type IContextWithCancel interface {
	context.Context
	Cancel()
	Canceled() bool
}

type contextWithCancel struct {
	context.Context
	cancel func()
}

func ContextWithCancel(parent context.Context) IContextWithCancel {
	ctx, cancel := context.WithCancel(parent)
	return &contextWithCancel {
		Context: ctx,
		cancel: cancel,
	}
}

func (c *contextWithCancel) Cancel() {
	c.cancel()
}

func (c *contextWithCancel) Canceled() bool {
	select {
	case <-c.Done():
		return true
	default:
		return false
	}
}


func ContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
