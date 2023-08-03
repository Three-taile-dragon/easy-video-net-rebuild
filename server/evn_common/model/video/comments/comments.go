package comments

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/video/barrage"
	"google.golang.org/appengine/user"
)

type Comment struct {
	common.PublicModel
	Uid            uint   `json:"uid" gorm:"column:uid"`
	VideoID        uint   `json:"video_id" gorm:"column:video_id"`
	Context        string `json:"context" gorm:"column:context"`
	CommentID      uint   `json:"comment_id" gorm:"column:comment_id"`
	CommentUserID  uint   `json:"comment_user_id" gorm:"column:comment_user_id"`
	CommentFirstID uint   `json:"comment_first_id" gorm:"column:comment_first_id"`

	UserInfo  user.User         `json:"user_info" gorm:"foreignKey:Uid"`
	VideoInfo barrage.VideoInfo `json:"video_info" gorm:"foreignKey:VideoID"`
}
type CommentList []Comment

func (Comment) TableName() string {
	return "lv_video_contribution_comments"
}
