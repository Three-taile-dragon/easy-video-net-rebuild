package database

type DbConn interface {
	Begin()
	Rollback() //事务回滚
	Commit()   //提交
}
