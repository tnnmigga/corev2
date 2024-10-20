package rdb

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/tnnmigga/corev2/conf"
)

const (
	RedisModeNormal      = "normal"
	RedisModeReplication = "replication"
	RedisModeSharding    = "sharding"

	MaxActiveConns = 256
)

func init() {
	e := conf.Scan("rdb", &map[string]any{})
	if e != nil {
		return
	}
	if err := initFromConf(); err != nil {
		panic(err)
	}
}

var conns = map[string]conn{}

type config struct {
	Addrs    []string
	Password string
	Index    int
	Mode     string
	Master   string
}

func initFromConf() error {
	data := conf.Map[config]("rdb", nil)
	for k, v := range data {
		db, err := newConn(v)
		if err != nil {
			return err
		}
		conns[k] = db
	}
	return nil
}

type IClient interface {
	redis.Cmdable
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error
	Process(ctx context.Context, cmd redis.Cmder) error
}

func newConn(c config) (conn, error) {
	var cli IClient
	switch c.Mode {
	case RedisModeNormal:
		cli = redis.NewClient(&redis.Options{
			Addr:     c.Addrs[0],
			Password: c.Password,
			DB:       c.Index,
		})
	case RedisModeReplication:
		cli = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:     c.Master,
			SentinelAddrs:  c.Addrs,
			Password:       c.Password,
			RouteByLatency: true,
		})
	case RedisModeSharding:
		cli = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          c.Addrs,
			Password:       c.Password,
			RouteByLatency: true,
		})
	default:
		panic("invalid mode")
	}
	return conn{cli: cli}, nil
}
