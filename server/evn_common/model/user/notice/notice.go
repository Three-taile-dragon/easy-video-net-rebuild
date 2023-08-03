package notice

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	"gorm.io/datatypes"
)

type Notice struct {
	common.PublicModel
	Uid     uint   `json:"uid" gorm:"column:uid"`
	Cid     uint   `json:"cid" gorm:"column:cid"`
	Type    string `json:"type" gorm:"column:type"`
	ToID    uint   `json:"to_id" gorm:"column:to_id"`
	ISRead  uint   `json:"is_read" gorm:"column:is_read"`
	Content string `json:"content" gorm:"column:content"`

	VideoInfo   VideoInfo `json:"videoInfo" gorm:"foreignKey:to_id"`
	UserInfo    user.User `json:"userinfo"  gorm:"foreignKey:cid"`
	ArticleInfo Article   `json:"articleInfo" gorm:"foreignKey:to_id"`
}

var (
	Online         = "online"         //上线时进行通知
	VideoComment   = "videoComment"   //视频评论
	VideoLike      = "videoLike"      //视频点赞
	ArticleComment = "articleComment" //文章评论
	ArticleLike    = "articleLike"    //文章点赞

)

type NoticesList []Notice

func (Notice) TableName() string {
	return "lv_users_notice"
}

// VideoInfo 临时加一个video模型解决依赖循环
type VideoInfo struct {
	common.PublicModel
	Uid   uint           `json:"uid" gorm:"uid"`
	Title string         `json:"title" gorm:"title"`
	Video datatypes.JSON `json:"video" gorm:"video"`
	Cover datatypes.JSON `json:"cover" gorm:"cover"`
}

func (VideoInfo) TableName() string {
	return "lv_video_contribution"
}

type Article struct {
	common.PublicModel
	Uid              uint           `json:"uid" gorm:"uid"`
	ClassificationID uint           `json:"classification_id"  gorm:"classification_id"`
	Title            string         `json:"title" gorm:"title"`
	Cover            datatypes.JSON `json:"cover" gorm:"cover"`
}

func (Article) TableName() string {
	return "lv_article_contribution"
}
