package video

import (
	"dragonsss.cn/evn_other/internal/data/barrage"
	"dragonsss.cn/evn_other/internal/data/comments"
	"dragonsss.cn/evn_other/internal/data/common"
	"dragonsss.cn/evn_other/internal/data/like"
	"dragonsss.cn/evn_other/internal/data/user"
	"gorm.io/datatypes"
)

type VideosContribution struct {
	common.PublicModel
	Uid           uint           `json:"uid" gorm:"column:uid"`
	Title         string         `json:"title" gorm:"column:title"`
	Video         datatypes.JSON `json:"video" gorm:"column:video"` //默认1080p
	Video720p     datatypes.JSON `json:"video_720p" gorm:"column:video_720p"`
	Video480p     datatypes.JSON `json:"video_480p" gorm:"column:video_480p"`
	Video360p     datatypes.JSON `json:"video_360p" gorm:"column:video_360p"`
	MediaID       string         `json:"media_id" gorm:"column:media_id"`
	Cover         datatypes.JSON `json:"cover" gorm:"column:cover"`
	VideoDuration int64          `json:"video_duration" gorm:"column:video_duration"`
	Reprinted     int8           `json:"reprinted" gorm:"column:reprinted"`
	Label         string         `json:"label" gorm:"column:label"`
	Introduce     string         `json:"introduce" gorm:"column:introduce"`
	Heat          int            `json:"heat" gorm:"column:heat"`

	UserInfo user.User            `json:"user_info" gorm:"foreignKey:Uid"`
	Likes    like.LikesList       `json:"likes" gorm:"foreignKey:VideoID" `
	Comments comments.CommentList `json:"comments" gorm:"foreignKey:VideoID"`
	Barrage  barrage.BarragesList `json:"barrage" gorm:"foreignKey:VideoID"`
}

type VideosContributionList []VideosContribution

func NewVideosContribution() *[]VideosContribution {
	return &[]VideosContribution{}
}

func (VideosContribution) TableName() string {
	return "lv_video_contribution"
}
