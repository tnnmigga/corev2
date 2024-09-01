package rdb

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/tnnmigga/corev2/conf"
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

var dbs = map[string]*redis.Pool{}

type config struct {
	Pass  string
	Host  string
	Port  int
	Index int
}

func initFromConf() error {
	data := conf.Map[config]("natsmq", nil)
	for k, v := range data {
		db, err := newPool(v)
		if err != nil {
			return err
		}
		dbs[k] = db
	}
	return nil
}

func newPool(c config) (*redis.Pool, error) {
	p := &redis.Pool{
		MaxIdle:     150,
		MaxActive:   100,
		IdleTimeout: 30 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), redis.DialPassword(c.Pass), redis.DialDatabase(c.Index))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	conn := p.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("PING"))
	if err != nil {
		return nil, err
	}
	if r != "PONG" {
		return nil, fmt.Errorf("redis %#v pong error %v", c, r)
	}
	return p, nil
}
