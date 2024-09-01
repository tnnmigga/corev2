package mdb

import (
	"fmt"
	"runtime/debug"

	"github.com/tnnmigga/corev2/conf"
	"github.com/tnnmigga/corev2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	if err := initFromConf(); err != nil {
		panic(err)
	}
}

func Use(name string) *conn {
	return dbs[name]
}

func Default() *conn {
	return dbs["default"]
}

func RecoverWithRollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		log.Errorf("panic %v, %s", r, debug.Stack())
		tx.Rollback()
	}
}

var dbs = map[string]*conn{}

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
	data := conf.Map[config]("mdb", nil)
	for k, v := range data {
		conn, err := newConn(v)
		if err != nil {
			panic(err)
		}
		dbs[k] = conn
	}
	return nil
}

func newConn(c config) (*conn, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True", c.User, c.Pass, c.Host, c.Port, c.Name)
	dialector := mysql.Open(dsn)
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
