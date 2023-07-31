package liveInfo

import (
	"dragonsss.cn/evn_other/internal/data/common"
	"gorm.io/datatypes"
)

type LiveInfo struct {
	common.PublicModel
	Uid   uint           `json:"uid" gorm:"column:uid"`
	Title string         `json:"title" gorm:"column:title"`
	Img   datatypes.JSON `json:"img" gorm:"type:json;comment:img"`
}

func (LiveInfo) TableName() string {
	return "lv_live_info"
}
