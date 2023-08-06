package repo

import (
	"context"
	"dragonsss.cn/evn_common/model/article"
	"dragonsss.cn/evn_common/model/liveInfo"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/attention"
	"dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_user/internal/database"
)

type UserRepo interface {
	SaveUser(conn database.DbConn, ctx context.Context, mem *user.User) error
	IsExistByEmail(ctx context.Context, email string) (bool, error)
	IsExistByNameAndEmail(ctx context.Context, name string) (bool, error)
	IsExistByName(ctx context.Context, name string) (bool, error)
	IsExistByMobile(ctx context.Context, mobile string) (bool, error)
	IsExistByID(ctx context.Context, id uint32) (bool, error)
	IsAttention(ctx context.Context, uid uint32, id uint32) (bool, error)
	IsExistUserLiveInfo(ctx context.Context, liveinfo *liveInfo.LiveInfo) (bool, error)
	CheckPassword(ctx context.Context, name string) (mem *user.User, err error)
	FindUserById(ctx context.Context, id int64) (*user.User, error)
	FindUserByName(ctx context.Context, name string) (*user.User, error)
	FindUserByEmail(ctx context.Context, email string) (*user.User, error)
	GetAttentionNum(ctx context.Context, id uint32) (*int64, error)
	GetVermicelliNum(ctx context.Context, id uint32) (*int64, error)
	GetVideoListBySpace(ctx context.Context, id uint32) (*video.VideosContributionList, error)
	GetArticleBySpace(ctx context.Context, id uint32) (*article.ArticlesContributionList, error)
	GetAttentionList(ctx context.Context, id uint32) (*attention.AttentionsList, error)
	GetAttentionListByIdArr(ctx context.Context, id uint32) (arr []uint, err error)
	GetVermicelliList(ctx context.Context, id uint32) (*attention.AttentionsList, error)
	UpdatePureZero(ctx context.Context, id int64, update map[string]interface{}) (bool, error)
	UpdateLoginTime(conn database.DbConn, ctx context.Context, name string) error
	UpdateUser(conn database.DbConn, ctx context.Context, mem *user.User) error
	UpdateUserAvatar(ctx context.Context, tmpUser *user.User) (bool, error)
	FindUserLiveInfo(ctx context.Context, id int64) (*liveInfo.LiveInfo, error)
	SaveUserLiveInfo(ctx context.Context, liveinfo *liveInfo.LiveInfo) (bool, error)
	UpdateUserLiveInfo(ctx context.Context, liveinfo *liveInfo.LiveInfo) (bool, error)
	Attention(ctx context.Context, aid uint32, uid uint32) (*attention.Attention, error)
}

type MemberRepo interface {
	SaveMember(ctx context.Context, member *user.User) error
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	FindMemberById(ctx context.Context, id int64) (*user.User, error)
	UpdateLoginTime(ctx context.Context, id int64) error
}
