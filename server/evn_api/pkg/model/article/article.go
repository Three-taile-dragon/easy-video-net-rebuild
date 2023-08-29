package article

import (
	"dragonsss.cn/evn_common/model/common"
	"time"
)

type GetArticleContributionListReceiveStruct struct {
	PageInfo common.PageInfo `json:"page_info"`
}

type GetArticleContributionListByUserResponseStruct struct {
	Id             uint      `json:"id"`
	Uid            uint      `json:"uid" `
	Title          string    `json:"title" `
	Cover          string    `json:"cover" `
	Label          []string  `json:"label" `
	Content        string    `json:"content"`
	IsComments     bool      `json:"is_comments"`
	Heat           int       `json:"heat"`
	LikesNumber    int       `json:"likes_number"`
	CommentsNumber int       `json:"comments_number"`
	Classification string    `json:"classification"`
	CreatedAt      time.Time `json:"created_at"`
}

type GetArticleContributionListByUserResponseList []GetArticleContributionListByUserResponseStruct

type GetArticleContributionListByUserReceiveStruct struct {
	UserID uint `json:"userID" binding:"required"`
}

type GetArticleCommentReceiveStruct struct {
	PageInfo  common.PageInfo `json:"pageInfo"`
	ArticleID uint            `json:"articleID" binding:"required"`
}

// 评论信息
type commentsInfo struct {
	ID              uint             `json:"id"`
	CommentID       uint             `json:"comment_id"`
	CommentFirstID  uint             `json:"comment_first_id"`
	CreatedAt       time.Time        `json:"created_at"`
	Context         string           `json:"context"`
	Uid             uint             `json:"uid"`
	Username        string           `json:"username"`
	Photo           string           `json:"photo"`
	CommentUserID   uint             `json:"comment_user_id"`
	CommentUserName string           `json:"comment_user_name"`
	LowerComments   commentsInfoList `json:"lowerComments"`
}

type commentsInfoList []*commentsInfo

type GetArticleContributionCommentsResponseStruct struct {
	Id             uint             `json:"id"`
	Comments       commentsInfoList `json:"comments"`
	CommentsNumber int              `json:"comments_number"`
}

// ArticleClassificationInfo 文章分类信息
type ArticleClassificationInfo struct {
	ID       uint                          `json:"id"`
	AID      uint                          `json:"aid"`
	Label    string                        `json:"label"`
	Children ArticleClassificationInfoList `json:"children"`
}

type ArticleClassificationInfoList []*ArticleClassificationInfo

type GetArticleTotalInfoResponseStruct struct {
	Classification    ArticleClassificationInfoList `json:"classification"`
	ArticleNum        int64                         `json:"article_num"`
	ClassificationNum int64                         `json:"classification_num"`
}

type GetArticleContributionByIDReceiveStruct struct {
	ArticleID uint `json:"articleID" binding:"required"`
}

type GetArticleContributionByIDResponseStruct struct {
	Id             uint             `json:"id"`
	Uid            uint             `json:"uid" `
	Title          string           `json:"title" `
	Cover          string           `json:"cover" `
	Label          []string         `json:"label" `
	Content        string           `json:"content"`
	IsComments     bool             `json:"is_comments"`
	Heat           int              `json:"heat"`
	LikesNumber    int              `json:"likes_number"`
	Comments       commentsInfoList `json:"comments"`
	CommentsNumber int              `json:"comments_number"`
	CreatedAt      time.Time        `json:"created_at"`
}

type CreateArticleContributionReceiveStruct struct {
	Cover                         string   `json:"cover" binding:"required"`
	CoverUploadType               string   `json:"coverUploadType" binding:"required"`
	ArticleContributionUploadType string   `json:"articleContributionUploadType" binding:"required"`
	Title                         string   `json:"title" binding:"required"`
	Label                         []string `json:"label" binding:"required"`
	Content                       string   `json:"content" binding:"required"`
	Comments                      *bool    `json:"comments"  binding:"required"`
	ClassificationID              uint     `json:"classification_id"`
}

type UpdateArticleContributionReceiveStruct struct {
	ID                            uint     `json:"id" binding:"required"`
	Cover                         string   `json:"cover" binding:"required"`
	CoverUploadType               string   `json:"coverUploadType" binding:"required"`
	ArticleContributionUploadType string   `json:"articleContributionUploadType" binding:"required"`
	Title                         string   `json:"title" binding:"required"`
	Label                         []string `json:"label" binding:"required"`
	Content                       string   `json:"content" binding:"required"`
	Comments                      *bool    `json:"comments"  binding:"required"`
	ClassificationID              uint     `json:"classification_id"`
}

type DeleteArticleByIDReceiveStruct struct {
	ID uint `json:"id"`
}

type ArticlesPostCommentReceiveStruct struct {
	ArticleID uint   `json:"article_id"`
	Content   string `json:"content"`
	ContentID uint   `json:"content_id"`
}

type GetArticleManagementListReceiveStruct struct {
	PageInfo common.PageInfo `json:"page_info"`
}

type GetArticleManagementListItem struct {
	ID               uint     `json:"id"`
	ClassificationID uint     `json:"classification_id"`
	Title            string   `json:"title"`
	Cover            string   `json:"cover"`
	CoverUrl         string   `json:"cover_url"`
	CoverUploadType  string   `json:"cover_upload_type"`
	Label            []string `json:"label"`
	Content          string   `json:"content"`
	IsComments       bool     `json:"is_comments" `
	Heat             int      `json:"heat"`
}

type GetArticleManagementListResponseStruct []GetArticleManagementListItem
