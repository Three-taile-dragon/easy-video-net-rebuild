package favorites

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/collect"
	"gorm.io/datatypes"
)

type Favorites struct {
	common.PublicModel
	Uid     uint           `json:"uid" gorm:"column:uid"`
	Title   string         `json:"title" gorm:"column:title"`
	Content string         `json:"content" gorm:"column:content"`
	Cover   datatypes.JSON `json:"cover" gorm:"type:json;comment:cover"`
	Max     int            `json:"max" gorm:"column:max"`

	UserInfo    user.User            `json:"userInfo" gorm:"foreignKey:Uid"`
	CollectList collect.CollectsList `json:"collectList"  gorm:"foreignKey:FavoritesID"`
}

type FavoriteList []Favorites

func (Favorites) TableName() string {
	return "lv_users_favorites"
}
