package chatList

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	"time"
)

type ChatsListInfo struct {
	common.PublicModel
	Uid         uint      `json:"uid" gorm:"column:uid"`
	Tid         uint      `json:"tid"  gorm:"column:tid"`
	Unread      int       `json:"unread" gorm:"column:unread"`
	LastMessage string    `json:"last_message" gorm:"column:last_message"`
	LastAt      time.Time `json:"last_at" gorm:"column:last_at"`

	ToUserInfo user.User `json:"toUserInfo"  gorm:"foreignKey:tid"`
}

type ChatList []ChatsListInfo

func (ChatsListInfo) TableName() string {
	return "lv_users_chat_list"
}
