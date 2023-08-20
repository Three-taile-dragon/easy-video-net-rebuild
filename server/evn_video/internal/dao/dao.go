package dao

import (
	"dragonsss.com/evn_video/internal/database"
	"dragonsss.com/evn_video/internal/database/gorms"
)

//事务的具体实现

type TransactionTmpl struct {
	conn database.DbConn
}

func (t *TransactionTmpl) Action(f func(conn database.DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback() //事务回滚
		return err
	}
	t.conn.Commit() //事务提交
	return nil
}

func NewTransaction() *TransactionTmpl {
	return &TransactionTmpl{
		conn: gorms.NewTran(),
	}
}
