package classification

import (
	"dragonsss.cn/evn_common/model/common"
)

type Classification struct {
	common.PublicModel
	AID   uint   `json:"a_id" gorm:"column:a_id"`
	Label string `json:"label" gorm:"column:label"`
}

type ClassificationsList []Classification

func (Classification) TableName() string {
	return "lv_article_classification"
}
