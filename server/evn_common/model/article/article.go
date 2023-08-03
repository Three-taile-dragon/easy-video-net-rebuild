package article

import (
	"dragonsss.cn/evn_common/model/article/classification"
	"dragonsss.cn/evn_common/model/article/comments"
	"dragonsss.cn/evn_common/model/article/like"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	"gorm.io/datatypes"
)

type ArticlesContribution struct {
	common.PublicModel
	Uid                uint           `json:"uid" gorm:"column:uid"`
	ClassificationID   uint           `json:"classification_id"  gorm:"column:classification_id"`
	Title              string         `json:"title" gorm:"column:title"`
	Cover              datatypes.JSON `json:"cover" gorm:"column:cover"`
	Label              string         `json:"label" gorm:"column:label"`
	Content            string         `json:"content" gorm:"column:content"`
	ContentStorageType string         `json:"content_Storage_Type" gorm:"column:content_storage_type"`
	IsComments         int8           `json:"is_comments" gorm:"column:is_comments"`
	Heat               int            `json:"heat" gorm:"column:heat"`

	//光联表

	UserInfo       user.User                     `json:"user_info" gorm:"foreignKey:Uid"`
	Likes          like.LikesList                `json:"likes" gorm:"foreignKey:ArticleID" `
	Comments       comments.CommentList          `json:"comments" gorm:"foreignKey:ArticleID"`
	Classification classification.Classification `json:"classification"  gorm:"foreignKey:ClassificationID"`
}

type ArticlesContributionList []ArticlesContribution

func (ArticlesContribution) TableName() string {
	return "lv_article_contribution"
}
