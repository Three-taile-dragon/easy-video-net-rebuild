package barrage

import (
	"dragonsss.cn/evn_common/model/common"
	"gorm.io/datatypes"
	"os/user"
)

type Barrage struct {
	common.PublicModel
	Uid     uint    `json:"uid" gorm:"column:uid"`
	VideoID uint    `json:"video_id" gorm:"column:video_id"`
	Time    float64 `json:"time" gorm:"column:time"`
	Author  string  `json:"author" gorm:"column:author"`
	Type    uint    `json:"type" gorm:"column:type"`
	Text    string  `json:"text" gorm:"column:text"`
	Color   uint    `json:"color" gorm:"column:color"`

	UserInfo  user.User `json:"user_info" gorm:"foreignKey:Uid"`
	VideoInfo VideoInfo `json:"video_info" gorm:"foreignKey:VideoID"`
}

type BarragesList []Barrage

func (Barrage) TableName() string {
	return "lv_video_contribution_barrage"
}

// VideoInfo 临时加一个video模型解决依赖循环
type VideoInfo struct {
	common.PublicModel
	Uid   uint           `json:"uid" gorm:"uid"`
	Title string         `json:"title" gorm:"title"`
	Video datatypes.JSON `json:"video" gorm:"video"`
	Cover datatypes.JSON `json:"cover" gorm:"cover"`
}

func (VideoInfo) TableName() string {
	return "lv_video_contribution"
}
