package rotograph

import "dragonsss.cn/evn_common/model/common"
import "gorm.io/datatypes"

type Rotograph struct {
	common.PublicModel
	Title string         `json:"title" gorm:"column:title"`
	Cover datatypes.JSON `json:"cover" gorm:"column:cover"`
	Color string         `json:"color" gorm:"column:color" `
	Type  string         `json:"type" gorm:"column:type"`
	ToId  uint           `json:"to_id" gorm:"column:to_id"`
}
type List []Rotograph

func New() *[]Rotograph {
	return &[]Rotograph{}
}

func (Rotograph) TableName() string {
	return "lv_home_rotograph"
}
