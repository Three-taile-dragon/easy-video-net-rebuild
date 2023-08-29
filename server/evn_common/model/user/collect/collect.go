package collect

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/video"
)

type Collect struct {
	common.PublicModel
	Uid         uint `json:"uid" gorm:"column:uid"`
	FavoritesID uint `json:"favorites_id" gorm:"column:favorites_id"`
	VideoID     uint `json:"video_id" gorm:"column:video_id"`

	UserInfo  user.User                `json:"userInfo" gorm:"foreignKey:Uid"`
	VideoInfo video.VideosContribution `json:"videoInfo" gorm:"foreignKey:VideoID"`
}

type CollectsList []Collect

func (Collect) TableName() string {
	return "lv_users_collect"
}
