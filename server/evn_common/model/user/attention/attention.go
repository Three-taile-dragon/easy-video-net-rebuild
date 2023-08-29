package attention

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
)

type Attention struct {
	common.PublicModel
	Uid         uint `json:"uid" gorm:"column:uid"`
	AttentionID uint `json:"attention_id" gorm:"column:attention_id"`

	UserInfo          user.User `json:"user_info" gorm:"foreignKey:Uid"`
	AttentionUserInfo user.User `json:"attention_user_info" gorm:"foreignKey:AttentionID"`
}

type AttentionsList []Attention

func (Attention) TableName() string {
	return "lv_users_attention"
}
