package repo

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_user/internal/database"
)

type UserRepo interface {
	SaveUser(conn database.DbConn, ctx context.Context, mem *user.User) error
	GetUserByEmail(ctx context.Context, email string) (bool, error)
	GetUserByNameAndEmail(ctx context.Context, name string) (bool, error)
	GetUserByName(ctx context.Context, name string) (bool, error)
	GetUserByMobile(ctx context.Context, mobile string) (bool, error)
	CheckPassword(ctx context.Context, name string) (mem *user.User, err error)
	FindUserById(ctx context.Context, id int64) (*user.User, error)
	FindUserByName(ctx context.Context, name string) (*user.User, error)
	UpdateLoginTime(conn database.DbConn, ctx context.Context, name string) error
}

type MemberRepo interface {
	SaveMember(ctx context.Context, member *user.User) error
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	FindMemberById(ctx context.Context, id int64) (*user.User, error)
	UpdateLoginTime(ctx context.Context, id int64) error
}
