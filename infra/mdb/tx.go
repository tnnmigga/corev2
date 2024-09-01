package mdb

import (
	"database/sql"

	"gorm.io/gorm"
)

// // 读未提交 // 不建议用
// func (c *conn) TxReadUncommit() *gorm.DB {
// 	tx := c.Begin(&sql.TxOptions{Isolation: sql.LevelReadUncommitted})
// 	return tx
// }

// 读已提交
func (c *conn) TxReadCommit() *gorm.DB {
	tx := c.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})
	return tx
}

// 可重复读
// 从这个级别开始,SELECT ... FOR UPDATE下会根据WHERE语句添加间隙锁
func (c *conn) TxRepeatableRead() *gorm.DB {
	tx := c.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	return tx
}

// 串行化
func (c *conn) TxSerializable() *gorm.DB {
	tx := c.Begin(&sql.TxOptions{Isolation: sql.LevelSerializable})
	return tx
}
