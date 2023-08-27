package repo

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/favorites"
	video2 "dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_common/model/video/barrage"
	"dragonsss.cn/evn_common/model/video/comments"
	"dragonsss.cn/evn_common/model/video/like"
	"dragonsss.cn/evn_grpc/video"
	"dragonsss.com/evn_video/internal/database"
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
	DeleteVideoByID(ctx context.Context, req *video.CommonIDAndUIDRequest) (bool, error)
	GetCommentFirstIDByID(ctx context.Context, contentID uint32) (*comments.Comment, error)
	GetCommentUserIDByID(ctx context.Context, contentID uint32) (*comments.Comment, error)
	CreateComment(conn database.DbConn, ctx context.Context, comment *comments.Comment) error
	AddNotice(ctx context.Context, uid uint, cid uint, tid uint, tp string, c string) error
	GetVideoManagementList(ctx context.Context, req *video.GetVideoManagementListRequest) (*video2.VideosContributionList, error)
	GetLikeVideo(conn database.DbConn, ctx context.Context, req *video.CommonIDAndUIDRequest) (*like.Likes, error)
	DeleteLikeVideo(conn database.DbConn, ctx context.Context, req *video.CommonIDAndUIDRequest, li *like.Likes) error
	DeleteNotice(ctx context.Context, videoUid uint, uid uint32, videoID uint32, videoLike string) error
	LikeVideo(conn database.DbConn, ctx context.Context, likeVideo *like.Likes) error
}
