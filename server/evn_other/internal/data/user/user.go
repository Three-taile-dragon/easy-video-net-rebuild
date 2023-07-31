package user

import (
	"dragonsss.cn/evn_other/internal/data/common"
	"dragonsss.cn/evn_other/internal/data/liveInfo"
	"gorm.io/datatypes"
	"time"
)

// User 表结构体
type User struct {
	common.PublicModel
	Email     string         `json:"email" gorm:"column:email"`
	Username  string         `json:"username" gorm:"column:username"`
	Openid    string         `json:"openid" gorm:"column:openid"`
	Salt      string         `json:"salt" gorm:"column:salt"`
	Password  string         `json:"password" gorm:"column:password"`
	Photo     datatypes.JSON `json:"photo" gorm:"column:photo"`
	Gender    int8           `json:"gender" gorm:"column:gender"`
	BirthDate time.Time      `json:"birth_date" gorm:"column:birth_date"`
	IsVisible int8           `json:"is_visible" gorm:"column:is_visible"`
	Signature string         `json:"signature" gorm:"column:signature"`

	LiveInfo liveInfo.LiveInfo `json:"liveInfo" gorm:"foreignKey:Uid"`
}

type UserList []User

func (User) TableName() string {
	return "lv_users"
}
