package chdb

import (
	"fmt"

	"github.com/tnnmigga/corev2/conf"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
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
	*gorm.DB
}

func initFromConf() error {
	data := conf.Map[config]("chdb", nil)
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
	dsn := fmt.Sprintf("clickhouse://%s%s:%d/%s?dial_timeout=200ms&max_execution_time=60", account, c.Host, c.Port, c.Name)
	dialector := clickhouse.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger{},
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(50)
	return &conn{DB: db}, nil
}
