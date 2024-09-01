package rdb

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type conn struct {
	redis.Conn
}

func Default() conn {
	p, ok := dbs["default"]
	if !ok {
		panic("redis name not found")
	}
	c := p.Get()
	return conn{Conn: c}
}

func Use(name string) conn {
	p, ok := dbs[name]
	if !ok {
		panic("redis name not found")
	}
	c := p.Get()
	return conn{Conn: c}
}

func (c conn) Do(commandName string, args ...interface{}) Result {
	reply, err := c.Conn.Do(commandName, args...)
	return Result{reply, err}
}

func (c conn) DoOnce(commandName string, args ...interface{}) Result {
	defer c.Close()
	reply, err := c.Conn.Do(commandName, args...)
	return Result{reply, err}
}

func (c conn) SET(key string, value any, extargs ...any) Result {
	return c.DoOnce("SET", append(append(make([]interface{}, 0, len(extargs)+2), key, value), extargs...)...)
}

func (c conn) SETNX(key string, value any) Result {
	return c.DoOnce("SETNX", key, value)
}

func (c conn) EXPIRE(key string, duration time.Duration) Result {
	return c.DoOnce("EXPIRE", key, duration/time.Second)
}

func (c conn) GET(key string) Result {
	return c.DoOnce("GET", key)
}

func (c conn) MGET(keys ...any) Result {
	return c.DoOnce("MGET", keys...)
}

func (c conn) DEL(keys ...any) Result {
	return c.DoOnce("DEL", keys...)
}

func (c conn) HSET(key string, field any, value any) Result {
	return c.DoOnce("HSET", key, field, value)
}

func (c conn) HSETNX(key string, field any, value any) Result {
	return c.DoOnce("HSETNX", key, field, value)
}

func (c conn) HGET(key string, field any) Result {
	return c.DoOnce("HGET", key, field)
}

func (c conn) HDEL(keys ...any) Result {
	return c.DoOnce("HDEL", keys...)
}

func (c conn) HMSET(key string, args ...any) Result {
	args = append([]any{key}, args...)
	return c.DoOnce("HMSET", args...)
}

func (c conn) HMGET(key string, fields ...any) Result {
	return c.DoOnce("HMGET", append(append(make([]interface{}, 0, len(fields)+1), key), fields...)...)
}

func (c conn) HGETALL(key string) Result {
	return c.DoOnce("HGETALL", key)
}

func (c conn) HINCRBY(key string, field string, incr any) Result {
	return c.DoOnce("HINCRBY", key, field, incr)
}

func (c conn) HEXISTS(key string, field any) Result {
	return c.DoOnce("HEXISTS", key, field)
}

func (c conn) LPUSH(key string, values ...any) Result {
	return c.DoOnce("LPUSH", append(append(make([]interface{}, 0, len(values)+1), key), values...)...)
}

func (c conn) RPUSH(key string, values ...any) Result {
	return c.DoOnce("RPUSH", append(append(make([]interface{}, 0, len(values)+1), key), values...)...)
}

func (c conn) LPOP(key string) Result {
	return c.DoOnce("LPOP", key)
}

func (c conn) RPOP(key string) Result {
	return c.DoOnce("RPOP", key)
}

func (c conn) LLEN(key string) Result {
	return c.DoOnce("LLEN", key)
}

func (c conn) LRANGE(key string, start int32, stop int32) Result {
	return c.DoOnce("LRANGE", key, start, stop)
}

func (c conn) SADD(key string, members ...any) Result {
	return c.DoOnce("SADD", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
}

func (c conn) SMEMBERS(key string) Result {
	return c.DoOnce("SMEMBERS", key)
}

func (c conn) SREM(key string, members ...any) Result {
	return c.DoOnce("SREM", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
}

func (c conn) INCR(key string) Result {
	return c.DoOnce("INCR", key)
}

func (c conn) ZADD(key string, score any, member any) Result {
	return c.DoOnce("ZADD", key, score, member)
}

func (c conn) ZRANGE(key string, start int32, stop int32) Result {
	return c.DoOnce("ZRANGE", key, start, stop)
}

func (c conn) ZRANGE_WITHSCORES(key string, start int32, stop int32) Result {
	return c.DoOnce("ZRANGE", key, start, stop, "WITHSCORES")
}

func (c conn) ZRANGEBYSCORE(key string, min float64, max float64) Result {
	return c.DoOnce("ZRANGEBYSCORE", key, min, max)
}

func (c conn) ZREM(key string, members ...any) Result {
	return c.DoOnce("ZREM", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
}

func (c conn) ZREMRANGEBYSCORE(key string, min float64, max float64) Result {
	return c.DoOnce("ZREMRANGEBYSCORE", key, min, max)
}

func (c conn) ZREMRANGEBYRANK(key string, start int32, stop int32) Result {
	return c.DoOnce("ZREMRANGEBYRANK", key, start, stop)
}

func (c conn) ZCARD(key string) Result {
	return c.DoOnce("ZCARD", key)
}

func (c conn) ZSCORE(key string, member any) Result {
	return c.DoOnce("ZSCORE", key, member)
}

func (c conn) ZRANK(key string, member any) Result {
	return c.DoOnce("ZRANK", key, member)
}

func (c conn) ZREVRANK(key string, member any) Result {
	return c.DoOnce("ZREVRANK", key, member)
}

func (c conn) ZREVRANGE(key string, start int32, stop int32) Result {
	return c.DoOnce("ZREVRANGE", key, start, stop)
}

func (c conn) ZREVRANGEBYSCORE(key string, max float64, min float64) Result {
	return c.DoOnce("ZREVRANGEBYSCORE", key, max, min)
}

func (c conn) ZCOUNT(key string, min float64, max float64) Result {
	return c.DoOnce("ZCOUNT", key, min, max)
}

func (c conn) ZINCRBY(key string, increment any, member any) Result {
	return c.DoOnce("ZINCRBY", key, increment, member)
}

func (c conn) ZRANKBYSCORE(key string, member any) Result {
	return c.DoOnce("ZRANKBYSCORE", key, member)
}

func (c conn) ZREVRANKBYSCORE(key string, member any) Result {
	return c.DoOnce("ZREVRANKBYSCORE", key, member)
}

func (c conn) Multi() *Multi {
	m := &Multi{conn: c.Conn}
	m.err = c.Send("MULTI")
	return m
}

type Multi struct {
	conn  redis.Conn
	count int
	err   error
}

func (m *Multi) Add(commandName string, args ...interface{}) {
	if m.err != nil {
		return
	}
	m.err = m.conn.Send(commandName, args...)
	m.count++
}

func (m *Multi) Commit() ([]Result, error) {
	defer m.conn.Close()
	if m.err != nil {
		return nil, m.err
	}
	m.err = m.conn.Send("EXEC")
	if m.err != nil {
		return nil, m.err
	}
	m.err = m.conn.Flush()
	if m.err != nil {
		return nil, m.err
	}
	for i := 0; i < m.count+1; i++ {
		_, m.err = m.conn.Receive()
		if m.err != nil {
			return nil, m.err
		}
	}
	reply, err := m.conn.Receive()
	if err != nil {
		return nil, err
	}
	values, err := redis.Values(reply, err)
	if err != nil {
		return nil, err
	}
	result := make([]Result, 0, len(values))
	for _, value := range values {
		result = append(result, Result{
			value: value,
		})
	}
	return result, nil
}

func (m *Multi) GET(key string) {
	m.Add("GET", key)
}

func (c *Multi) MGET(keys ...any) {
	c.Add("MGET", keys...)
}

func (m *Multi) SET(key string, value any, extargs ...any) {
	m.Add("SET", append(append(make([]interface{}, 0, len(extargs)+2), key, value), extargs...)...)
}

func (m *Multi) SETNX(key string, value any) {
	m.Add("SETNX", key, value)
}

func (m *Multi) EXPIRE(key string, duration time.Duration) {
	m.Add("EXPIRE", key, duration/time.Second)
}

func (m *Multi) DEL(keys ...any) {
	m.Add("DEL", keys...)
}

func (m *Multi) HSET(key string, field any, value any) {
	m.Add("HSET", key, field, value)
}

func (m *Multi) HSETNX(key string, field any, value any) {
	m.Add("HSETNX", key, field, value)
}

func (m *Multi) HGET(key string, field any) {
	m.Add("HGET", key, field)
}

func (m *Multi) HDEL(keys ...any) {
	m.Add("HDEL", keys...)
}

func (m *Multi) HMSET(key string, args ...any) {
	args = append([]any{key}, args...)
	m.Add("HMSET", args...)
}

func (m *Multi) HMGET(key string, fields ...any) {
	m.Add("HMGET", append(append(make([]interface{}, 0, len(fields)+1), key), fields...)...)
}

func (m *Multi) HGETALL(key string) {
	m.Add("HGETALL", key)
}

func (m *Multi) HINCRBY(key string, field string, incr any) {
	m.Add("HINCRBY", key, field, incr)
}

func (m *Multi) HEXISTS(key string, field any) {
	m.Add("HEXISTS", key, field)
}

func (m *Multi) LPUSH(key string, values ...any) {
	m.Add("LPUSH", append(append(make([]interface{}, 0, len(values)+1), key), values...)...)
}

func (m *Multi) RPUSH(key string, values ...any) {
	m.Add("RPUSH", append(append(make([]interface{}, 0, len(values)+1), key), values...)...)
}

func (m *Multi) LPOP(key string) {
	m.Add("LPOP", key)
}

func (m *Multi) RPOP(key string) {
	m.Add("RPOP", key)
}

func (m *Multi) LLEN(key string) {
	m.Add("LLEN", key)
}

func (m *Multi) LRANGE(key string, start int32, stop int32) {
	m.Add("LRANGE", key, start, stop)
}

func (m *Multi) SADD(key string, members ...any) {
	m.Add("SADD", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
}

func (m *Multi) SMEMBERS(key string) {
	m.Add("SMEMBERS", key)
}

func (m *Multi) SREM(key string, members ...any) {
	m.Add("SREM", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
}

func (m *Multi) INCR(key string) {
	m.Add("INCR", key)
}

func (m *Multi) ZADD(key string, score float64, member any) {
	m.Add("ZADD", key, score, member)
}

func (m *Multi) ZRANGE(key string, start int32, stop int32) {
	m.Add("ZRANGE", key, start, stop)
}

func (m *Multi) ZRANGE_WITHSCORES(key string, start int32, stop int32) {
	m.Add("ZRANGE", key, start, stop, "WITHSCORES")
}
func (m *Multi) ZRANGEBYSCORE(key string, min float64, max float64) {
	m.Add("ZRANGEBYSCORE", key, min, max)
}

func (m *Multi) ZREM(key string, members ...any) {
	m.Add("ZREM", append(append(make([]interface{}, 0, len(members)+1), key), members...)...)
}

func (m *Multi) ZREMRANGEBYSCORE(key string, min float64, max float64) {
	m.Add("ZREMRANGEBYSCORE", key, min, max)
}

func (m *Multi) ZREMRANGEBYRANK(key string, start int32, stop int32) {
	m.Add("ZREMRANGEBYRANK", key, start, stop)
}

func (m *Multi) ZCARD(key string) {
	m.Add("ZCARD", key)
}

func (m *Multi) ZSCORE(key string, member any) {
	m.Add("ZSCORE", key, member)
}

func (m *Multi) ZRANK(key string, member any) {
	m.Add("ZRANK", key, member)
}

func (m *Multi) ZREVRANK(key string, member any) {
	m.Add("ZREVRANK", key, member)
}

func (m *Multi) ZREVRANGE(key string, start int32, stop int32) {
	m.Add("ZREVRANGE", key, start, stop)
}

func (m *Multi) ZREVRANGEBYSCORE(key string, max float64, min float64) {
	m.Add("ZREVRANGEBYSCORE", key, max, min)
}

func (m *Multi) ZCOUNT(key string, min float64, max float64) {
	m.Add("ZCOUNT", key, min, max)
}

func (m *Multi) ZINCRBY(key string, increment float64, member any) {
	m.Add("ZINCRBY", key, increment, member)
}

func (m *Multi) ZRANKBYSCORE(key string, member any) {
	m.Add("ZRANKBYSCORE", key, member)
}

func (m *Multi) ZREVRANKBYSCORE(key string, member any) {
	m.Add("ZREVRANKBYSCORE", key, member)
}
