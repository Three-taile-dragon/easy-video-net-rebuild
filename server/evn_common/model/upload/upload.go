package upload

import "dragonsss.cn/evn_common/model/common"

type Upload struct {
	common.PublicModel
	Interfaces string  `json:"interface"  gorm:"column:interface"`
	Method     string  `json:"method"  gorm:"column:method"`
	Path       string  `json:"path" gorm:"column:path"`
	Quality    float64 `json:"quality"  gorm:"column:quality"`
}

func (Upload) TableName() string {
	return "lv_upload_method"
}
