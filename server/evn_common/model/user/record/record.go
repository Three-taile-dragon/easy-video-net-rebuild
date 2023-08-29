package record

import (
	"dragonsss.cn/evn_common/model/article"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/video"
)

type Record struct {
	common.PublicModel
	Uid  uint   `json:"column:uid"`
	Type string `json:"type" gorm:"column:type"`
	ToId uint   `json:"to_id" gorm:"column:to_id"`

	VideoInfo   video.VideosContribution     `json:"videoInfo" gorm:"foreignKey:to_id"`
	Userinfo    user.User                    `json:"users.User"  gorm:"foreignKey:uid"`
	ArticleInfo article.ArticlesContribution `json:"articleInfo" gorm:"foreignKey:to_id"`
}

type RecordsList []Record

func (Record) TableName() string {
	return "lv_users_record"
}
