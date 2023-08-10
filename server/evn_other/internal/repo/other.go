package repo

import (
	"context"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/rotograph"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_grpc/other"
)

type OtherRepo interface {
	FindRotographInfo(ctx context.Context) (*rotograph.List, error)
	FindVideoList(ctx context.Context, info common.PageInfo) (*video.VideosContributionList, error)
	FindLiveInfo(ctx context.Context, req *other.CommonIDAndUIDRequest) (*user.User, error)
	AddLiveRecord(ctx context.Context, uid uint32, id uint32) error
	GetBeLiveList(ctx context.Context, keys []uint) (*user.UserList, error)
}
