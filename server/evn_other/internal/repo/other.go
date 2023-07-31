package repo

import (
	"context"
	"dragonsss.cn/evn_other/internal/data/common"
	"dragonsss.cn/evn_other/internal/data/rotograph"
	"dragonsss.cn/evn_other/internal/data/video"
)

type OtherRepo interface {
	FindRotographInfo(ctx context.Context) (*rotograph.List, error)
	FindVideoList(ctx context.Context, info common.PageInfo) (*video.VideosContributionList, error)
}
