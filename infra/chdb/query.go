package chdb

import (
	"runtime/debug"

	"github.com/tnnmigga/corev2/log"
	"gorm.io/gorm"
)

func Use(name string) *conn {
	return conns[name]
}

func Default() *conn {
	return conns["default"]
}

func RecoverWithRollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		log.Errorf("panic %v, %s", r, debug.Stack())
		tx.Rollback()
	}
}
