package mysql

import (
	"context"
	"dragonsss.cn/evn_user/internal/data/user"
	"dragonsss.cn/evn_user/internal/database"
	"dragonsss.cn/evn_user/internal/database/gorms"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	conn *gorms.GormConn
}

func NewUserDao() *UserDao {
	return &UserDao{
		conn: gorms.New(),
	}
}

func (u *UserDao) SaveUser(conn database.DbConn, ctx context.Context, mem *user.User) error {
	u.conn = conn.(*gorms.GormConn) //使用事务操作
	return u.conn.Tx(ctx).Create(mem).Error
}

func (u *UserDao) GetUserByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("email=?", email).Count(&count).Error //数据库查询
	return count > 0, err
}

func (u *UserDao) GetUserByAccountAndEmail(ctx context.Context, account string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("email=? or account=?", account, account).Count(&count).Error //数据库查询
	return count > 0, err
}

func (u *UserDao) GetUserByAccount(ctx context.Context, account string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("account=?", account).Count(&count).Error //数据库查询
	return count > 0, err
}

func (u *UserDao) GetUserByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("name=?", name).Count(&count).Error //数据库查询
	return count > 0, err
}

func (u *UserDao) GetUserByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("mobile=?", mobile).Count(&count).Error //数据库查询
	return count > 0, err
}
func (u *UserDao) FindUser(ctx context.Context, account string, pwd string) (*user.User, error) {
	var mem *user.User
	err := u.conn.Session(ctx).Where("account=? and password=?", account, pwd).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return nil, nil
	}
	return mem, err
}

func (u *UserDao) FindUserById(ctx context.Context, id int64) (*user.User, error) {
	var mem *user.User
	err := u.conn.Session(ctx).Where("id=?", id).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return nil, nil
	}
	return mem, err
}

func (u *UserDao) UpdateLoginTime(conn database.DbConn, ctx context.Context, id int64) error {
	u.conn = conn.(*gorms.GormConn) //使用事务操作
	mem, err := u.FindUserById(ctx, id)
	if err != nil {
		return err
	}
	mem.LastLoginTime = time.Now().UnixMilli()
	err = u.SaveUser(conn, ctx, mem)
	return err
}
