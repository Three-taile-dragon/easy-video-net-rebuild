package repo

import (
	"context"
	"dragonsss.cn/evn_common/model/article"
	comments2 "dragonsss.cn/evn_common/model/article/comments"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/rotograph"
	"dragonsss.cn/evn_common/model/upload"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/attention"
	"dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_common/model/video/barrage"
	"dragonsss.cn/evn_common/model/video/comments"
	"dragonsss.cn/evn_grpc/other"
)

type OtherRepo interface {
	FindRotographInfo(ctx context.Context) (*rotograph.List, error)
	FindVideoList(ctx context.Context, info common.PageInfo) (*video.VideosContributionList, error)
	FindLiveInfo(ctx context.Context, req *other.CommonIDAndUIDRequest) (*user.User, error)
	AddLiveRecord(ctx context.Context, uid uint32, id uint32) error
	GetBeLiveList(ctx context.Context, keys []uint) (*user.UserList, error)
	FindDiscussVideoCommentList(ctx context.Context, uid uint32) (*video.VideosContributionList, error)
	GetVideoCommentListByIDs(ctx context.Context, videoIDs []uint, req *other.CommonDiscussRequest) (*comments.CommentList, error)
	FindDiscussArticleCommentList(ctx context.Context, uid uint32) (*article.ArticlesContributionList, error)
	GetArticleCommentListByIDs(ctx context.Context, articleIDs []uint, req *other.CommonDiscussRequest) (*comments2.CommentList, error)
	GetVideoBarrageListByIDs(ctx context.Context, videoIDs []uint, req *other.CommonDiscussRequest) (*barrage.BarragesList, error)
	FindUploadMethod(ctx context.Context, method string) (*upload.Upload, error)
	SearchVideo(ctx context.Context, page int32, size int32, keyword string) (*video.VideosContributionList, error)
	SearchUser(ctx context.Context, page int32, size int32, keyword string) (*user.UserList, error)
	GetAttentionList(ctx context.Context, uid uint32) (*attention.AttentionsList, error)
}
