package rdb

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var defaultCtx = context.Background()

type conn struct {
	cli IClient
}

func Default() *Cmd {
	db, ok := conns["default"]
	if !ok {
		panic("redis name not found")
	}
	return db.Cmd()
}

func Use(name string) *Cmd {
	db, ok := conns[name]
	if !ok {
		panic("redis name not found")
	}
	return db.Cmd()
}

type Cmd struct {
	baseCmd
}

func (c *Cmd) Origin() IClient {
	return c.cli
}

func (c conn) Cmd() *Cmd {
	cmder := &Cmd{baseCmd: baseCmd{ctx: defaultCtx}}
	cmder.do = func(op string, args ...any) {
		cmder.result = []*redis.Cmd{redis.NewCmd(cmder.ctx, append([]any{op}, args...))}
		_ = c.cli.Process(cmder.ctx, cmder.result[0])
	}
	return cmder
}

type Pipeliner struct {
	baseCmd
	cmds [][]any
}

func (c *Cmd) Pipeline() *Pipeliner {
	cmder := &Pipeliner{baseCmd: c.baseCmd}
	cmder.do = func(op string, args ...any) {
		cmder.cmds = append(cmder.cmds, append([]any{op}, args...))
	}
	cmder.process = func() {
		pipe := c.cli.Pipeline()
		for _, cmd := range cmder.cmds {
			pipe.Do(cmder.ctx, cmd...)
		}
		result, err := pipe.Exec(cmder.ctx)
		if err != nil {
			cmder.err = err
		}
		for _, item := range result {
			cmder.result = append(cmder.result, item.(*redis.Cmd))
		}
	}
	return cmder
}

type Mulit struct {
	Pipeliner
}

func (c *Cmd) Multi() *Mulit {
	cmder := &Mulit{Pipeliner: Pipeliner{baseCmd: c.baseCmd}}
	cmder.do = func(op string, args ...any) {
		cmder.cmds = append(cmder.cmds, append([]any{op}, args...))
	}
	cmder.process = func() {
		pipe := c.cli.TxPipeline()
		for _, cmd := range cmder.cmds {
			pipe.Do(cmder.ctx, cmd...)
		}
		result, err := pipe.Exec(cmder.ctx)
		if err != nil {
			cmder.err = err
		}
		for _, item := range result {
			cmder.result = append(cmder.result, item.(*redis.Cmd))
		}
	}
	return cmder
}

type Tx struct {
	Pipeliner
	watch []string
	txDo  func(*Tx) error
	retry int
	wait  time.Duration
}

func (c *Cmd) Tx(do func(tx *Tx) error) *Tx {
	cmder := &Tx{Pipeliner: Pipeliner{baseCmd: c.baseCmd}, txDo: do}
	cmder.do = func(op string, args ...any) {
		cmder.cmds = append(cmder.cmds, append([]any{op}, args...))
	}
	cmder.process = func() {
		txf := func(tx *redis.Tx) error {
			cmder.cmds = cmder.cmds[:0]
			err := cmder.txDo(cmder)
			if err != nil {
				return err
			}
			pipe := tx.TxPipeline()
			for _, cmd := range cmder.cmds {
				pipe.Do(cmder.ctx, cmd...)
			}
			result, err := pipe.Exec(cmder.ctx)
			if err != nil {
				return err
			}
			for _, item := range result {
				cmder.result = append(cmder.result, item.(*redis.Cmd))
			}
			return nil
		}

		for i := 0; i <= cmder.retry; i++ {
			cmder.err = c.cli.Watch(cmder.ctx, txf, cmder.watch...)
			if cmder.err == nil {
				// Success.
				return
			}
			if cmder.err != redis.TxFailedErr {
				return
			}
			if cmder.wait == 0 {
				time.Sleep(cmder.wait)
			}
		}
	}
	return cmder
}

func (tx *Tx) Watch(keys ...string) *Tx {
	tx.watch = keys
	return tx
}

func (tx *Tx) Retry(count int, wait time.Duration) *Tx {
	tx.retry = count
	tx.wait = wait
	return tx
}

type baseCmd struct {
	ctx     context.Context
	cli     IClient
	do      func(op string, args ...any)
	process func()
	result  []*redis.Cmd
	err     error
}

func (c *baseCmd) Exec() *baseCmd {
	if c.process != nil {
		c.process()
	}
	return c
}

func (c *baseCmd) WithContext(ctx context.Context) *baseCmd {
	c.ctx = ctx
	return c
}

func (c *baseCmd) Value() *redis.Cmd {
	return c.result[0]
}

func (c *baseCmd) Values() ([]*redis.Cmd, error) {
	return c.result, c.err
}

func (c *baseCmd) Error() error {
	if c.err != nil {
		return c.err
	}
	for _, cmd := range c.result {
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}

func (c *baseCmd) SET(key string, value any, extargs ...any) *baseCmd {
	c.do("SET", append(append(make([]interface{}, 0, len(extargs)+2), key, value), extargs...)...)
	return c
}

func (c *baseCmd) SETNX(key string, value any) *baseCmd {
	c.do("SETNX", key, value)
	return c
}

func (c *baseCmd) EXPIRE(key string, duration time.Duration) *baseCmd {
	c.do("EXPIRE", key, duration/time.Second)
	return c
}

func (c *baseCmd) GET(key string) *baseCmd {
	c.do("GET", key)
	return c
}

func (c *baseCmd) MGET(keys ...any) *baseCmd {
	c.do("MGET", keys...)
	return c
}

func (c *baseCmd) DEL(keys ...any) *baseCmd {
	c.do("DEL", keys...)
	return c
}

func (c *baseCmd) HSET(key string, field any, value any) *baseCmd {
	c.do("HSET", key, field, value)
	return c
}

func (c *baseCmd) HSETNX(key string, field any, value any) *baseCmd {
	c.do("HSETNX", key, field, value)
	return c
}

func (c *baseCmd) HGET(key string, field any) *baseCmd {
	c.do("HGET", key, field)
	return c
}

func (c *baseCmd) HDEL(keys ...any) *baseCmd {
	c.do("HDEL", keys...)
	return c
}

func (c *baseCmd) HMSET(key string, args ...any) *baseCmd {
	args = append([]any{key}, args...)
	c.do("HMSET", args...)
	return c
}

func (c *baseCmd) HMGET(key string, fields ...any) *baseCmd {
	c.do("HMGET", append(append(make([]interface{}, 0, len(fields)+1), key), fields...)...)
	return c
}

func (c *baseCmd) HGETALL(key string) *baseCmd {
	c.do("HGETALL", key)
	return c
}

func (c *baseCmd) HINCRBY(key string, field string, incr any) *baseCmd {
	c.do("HINCRBY", key, field, incr)
	return c
}

func (c *baseCmd) HEXISTS(key string, field any) *baseCmd {
	c.do("HEXISTS", key, field)
	return c
}

func (c *baseCmd) LPUSH(key string, values ...any) *baseCmd {
	c.do("LPUSH", append(append(make([]interface{}, 0, len(values)+1), key), values...)...)
	return c
}

func (c *baseCmd) RPUSH(key string, values ...any) *baseCmd {
	c.do("RPUSH", append(append(make([]interface{}, 0, len(values)+1), key), values...)...)
	return c
}

func (c *baseCmd) LPOP(key string) *baseCmd {
	c.do("LPOP", key)
	return c
}

func (c *baseCmd) RPOP(key string) *baseCmd {
	c.do("RPOP", key)
	return c
}

func (c *baseCmd) LLEN(key string) *baseCmd {
	c.do("LLEN", key)
	return c
}

func (c *baseCmd) LRANGE(key string, start int32, stop int32) *baseCmd {
	c.do("LRANGE", key, start, stop)
	return c
}

func (c *baseCmd) SADD(key string, members ...any) *baseCmd {
	c.do("SADD", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
	return c
}

func (c *baseCmd) SMEMBERS(key string) *baseCmd {
	c.do("SMEMBERS", key)
	return c
}

func (c *baseCmd) SREM(key string, members ...any) *baseCmd {
	c.do("SREM", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
	return c
}

func (c *baseCmd) INCR(key string) *baseCmd {
	c.do("INCR", key)
	return c
}

func (c *baseCmd) ZADD(key string, score any, member any) *baseCmd {
	c.do("ZADD", key, score, member)
	return c
}

func (c *baseCmd) ZRANGE(key string, start int32, stop int32) *baseCmd {
	c.do("ZRANGE", key, start, stop)
	return c
}

func (c *baseCmd) ZRANGE_WITHSCORES(key string, start int32, stop int32) *baseCmd {
	c.do("ZRANGE", key, start, stop, "WITHSCORES")
	return c
}

func (c *baseCmd) ZRANGEBYSCORE(key string, min float64, max float64) *baseCmd {
	c.do("ZRANGEBYSCORE", key, min, max)
	return c
}

func (c *baseCmd) ZREM(key string, members ...any) *baseCmd {
	c.do("ZREM", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
	return c
}

func (c *baseCmd) ZREMRANGEBYSCORE(key string, min float64, max float64) *baseCmd {
	c.do("ZREMRANGEBYSCORE", key, min, max)
	return c
}

func (c *baseCmd) ZREMRANGEBYRANK(key string, start int32, stop int32) *baseCmd {
	c.do("ZREMRANGEBYRANK", key, start, stop)
	return c
}

func (c *baseCmd) ZCARD(key string) *baseCmd {
	c.do("ZCARD", key)
	return c
}

func (c *baseCmd) ZSCORE(key string, member any) *baseCmd {
	c.do("ZSCORE", key, member)
	return c
}

func (c *baseCmd) ZRANK(key string, member any) *baseCmd {
	c.do("ZRANK", key, member)
	return c
}

func (c *baseCmd) ZREVRANK(key string, member any) *baseCmd {
	c.do("ZREVRANK", key, member)
	return c
}

func (c *baseCmd) ZREVRANGE(key string, start int32, stop int32) *baseCmd {
	c.do("ZREVRANGE", key, start, stop)
	return c
}

func (c *baseCmd) ZREVRANGEBYSCORE(key string, max float64, min float64) *baseCmd {
	c.do("ZREVRANGEBYSCORE", key, max, min)
	return c
}

func (c *baseCmd) ZCOUNT(key string, min float64, max float64) *baseCmd {
	c.do("ZCOUNT", key, min, max)
	return c
}

func (c *baseCmd) ZINCRBY(key string, increment any, member any) *baseCmd {
	c.do("ZINCRBY", key, increment, member)
	return c
}

func (c *baseCmd) ZRANKBYSCORE(key string, member any) *baseCmd {
	c.do("ZRANKBYSCORE", key, member)
	return c
}

func (c *baseCmd) ZREVRANKBYSCORE(key string, member any) *baseCmd {
	c.do("ZREVRANKBYSCORE", key, member)
	return c
}
