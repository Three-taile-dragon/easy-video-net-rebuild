package repo

import (
	"context"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/rotograph"
	"dragonsss.cn/evn_common/model/video"
)

type OtherRepo interface {
	FindRotographInfo(ctx context.Context) (*rotograph.List, error)
	FindVideoList(ctx context.Context, info common.PageInfo) (*video.VideosContributionList, error)
}
