package comments

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	"gorm.io/datatypes"
)

type Comment struct {
	common.PublicModel
	Uid            uint   `json:"uid" gorm:"column:uid"`
	ArticleID      uint   `json:"article_id" gorm:"column:article_id"`
	Context        string `json:"context" gorm:"column:context"`
	CommentID      uint   `json:"comment_id" gorm:"column:comment_id"`
	CommentUserID  uint   `json:"comment_user_id" gorm:"column:comment_user_id"`
	CommentFirstID uint   `json:"comment_first_id" gorm:"column:comment_first_id"`

	UserInfo    user.User `json:"user_info" gorm:"foreignKey:Uid"`
	ArticleInfo Article   `json:"article_info" gorm:"foreignKey:ArticleID"`
}

type CommentList []Comment

func (Comment) TableName() string {
	return "lv_article_contribution_comments"
}

type Article struct {
	common.PublicModel
	Uid              uint           `json:"uid" gorm:"uid"`
	ClassificationID uint           `json:"classification_id"  gorm:"classification_id"`
	Title            string         `json:"title" gorm:"title"`
	Cover            datatypes.JSON `json:"cover" gorm:"cover"`
}

func (Article) TableName() string {
	return "lv_article_contribution"
}
