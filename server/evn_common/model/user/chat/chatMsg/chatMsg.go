package chatMsg

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
)

type Msg struct {
	common.PublicModel
	Uid     uint   `json:"uid" gorm:"column:uid"`
	Tid     uint   `json:"tid"  gorm:"column:tid"`
	Type    string `json:"type" gorm:"column:type"`
	Message string `json:"message" gorm:"column:message"`

	UInfo user.User `json:"UInfo"  gorm:"foreignKey:uid"`
	TInfo user.User `json:"TInfo"  gorm:"foreignKey:tid"`
}

type MsgList []Msg

func (Msg) TableName() string {
	return "lv_users_chat_msg"
}
