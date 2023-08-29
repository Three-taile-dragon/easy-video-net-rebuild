package repo

import (
	"context"
	article2 "dragonsss.cn/evn_common/model/article"
	"dragonsss.cn/evn_common/model/article/classification"
	comments2 "dragonsss.cn/evn_common/model/article/comments"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_grpc/article"
	"dragonsss.com/evn_article/internal/database"
)

type ArticleRepo interface {
	GetUserByID(ctx context.Context, uid uint32) (*user.User, error)
	GetArticleList(ctx context.Context, req *article.CommonPageInfo) (*article2.ArticlesContributionList, error)
	GetArticleListByID(ctx context.Context, req *article.CommonIDRequest) (*article2.ArticlesContributionList, error)
	GetArticleComments(ctx context.Context, req *article.GetArticleCommentRequest) (*article2.ArticlesContribution, error)
	FindAllClassification(ctx context.Context, req *article.CommonZeroRequest) (*classification.ClassificationsList, error)
	GetAllArticleCount(ctx context.Context) (*int64, error)
	GetArticleInfoByID(ctx context.Context, id uint32) (*article2.ArticlesContribution, error)
	AddArticleRecord(ctx context.Context, req *article.CommonIDAndUIDRequest) (bool, error)
	WatchArticle(ctx context.Context, id uint32) error
	CreateArticle(ctx context.Context, articlesContribution *article2.ArticlesContribution) error
	UpdateArticle(ctx context.Context, m *map[string]interface{}) error
	DeleteArticleByID(ctx context.Context, req *article.CommonIDAndUIDRequest) (bool, error)
	GetCommentFirstIDByID(ctx context.Context, contentID uint32) (*comments2.Comment, error)
	GetCommentUserIDByID(ctx context.Context, contentID uint32) (*comments2.Comment, error)
	CreateArticleComment(conn database.DbConn, ctx context.Context, c *comments2.Comment) error
	AddNotice(ctx context.Context, uid uint, cid uint, tid uint, tp string, c string) error
	GetArticleManagementList(ctx context.Context, req *article.GetArticleManagementListRequest) (*article2.ArticlesContributionList, error)
}
