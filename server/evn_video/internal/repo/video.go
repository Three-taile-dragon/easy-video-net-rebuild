package repo

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/favorites"
	video2 "dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_common/model/video/barrage"
	"dragonsss.cn/evn_grpc/video"
)

type VideoRepo interface {
	GetVideoBarrageByID(ctx context.Context, id uint32) (*barrage.BarragesList, error)
	GetVideoComments(ctx context.Context, req *video.GetVideoCommentRequest) (*video2.VideosContribution, error)
	GetUserByID(ctx context.Context, uid uint32) (*user.User, error)
	FindVideoByID(ctx context.Context, videoID uint32) (*video2.VideosContribution, error)
	WatchVideo(ctx context.Context, videoID uint32) (bool, error)
	IsAttention(ctx context.Context, uid uint32, id uint) (bool, error)
	IsLike(ctx context.Context, uid uint32, id uint) (bool, error)
	GetFavoritesList(ctx context.Context, uid uint32) (*favorites.FavoriteList, error)
	FindIsCollectByFavorites(ctx context.Context, id uint32, flIDs []uint) (bool, error)
	AddVideoRecord(ctx context.Context, uid uint32, videoID uint32) error
	GetRecommendList(ctx context.Context) (*video2.VideosContributionList, error)
	CreateVideoBarrage(ctx context.Context, bg *barrage.Barrage) (bool, error)
	CreateVideo(ctx context.Context, videoContribution *video2.VideosContribution) (bool, error)
	UpdateVideo(ctx context.Context, videoContribution *video2.VideosContribution) (bool, error)
}
