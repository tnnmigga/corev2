package mgdb

import (
	"context"
	"fmt"
	"time"

	"github.com/tnnmigga/corev2/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	if err := initFromConf(); err != nil {
		panic(err)
	}
}

var conns = map[string]*conn{}

type config struct {
	User string
	Pass string
	Host string
	Port int
	Name string
}

type conn struct {
	*mongo.Database
}

func Use(name string) *conn {
	return conns[name]
}

func Default() *conn {
	return conns["default"]
}

func initFromConf() error {
	data := conf.Map[config]("mgdb", nil)
	for k, v := range data {
		conn, err := newConn(v)
		if err != nil {
			panic(err)
		}
		conns[k] = conn
	}
	return nil
}

func newConn(c config) (*conn, error) {
	var account string
	if len(c.User) > 0 {
		account = fmt.Sprintf("%s:%s@", c.User, c.Pass)
	}
	uri := fmt.Sprintf("mongodb://%s%s:%d/%s", account, c.Host, c.Port, c.Name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := cli.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	database := cli.Database(c.Name)
	return &conn{Database: database}, nil
}
