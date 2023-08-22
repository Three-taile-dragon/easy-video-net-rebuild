package video

import (
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	"gorm.io/datatypes"
	"time"
)

type GetVideoBarrageReceiveStruct struct {
	ID string `json:"id"`
}

type Barrage struct {
	common.PublicModel
	Uid     uint    `json:"uid" gorm:"column:uid"`
	VideoID uint    `json:"video_id" gorm:"column:video_id"`
	Time    float64 `json:"time" gorm:"column:time"`
	Author  string  `json:"author" gorm:"column:author"`
	Type    uint    `json:"type" gorm:"column:type"`
	Text    string  `json:"text" gorm:"column:text"`
	Color   uint    `json:"color" gorm:"column:color"`

	UserInfo  user.User `json:"user_info" gorm:"foreignKey:Uid"`
	VideoInfo VideoInfo `json:"video_info" gorm:"foreignKey:VideoID"`
}

// VideoInfo 临时加一个video模型解决依赖循环
type VideoInfo struct {
	common.PublicModel
	Uid   uint           `json:"uid" gorm:"uid"`
	Title string         `json:"title" gorm:"title"`
	Video datatypes.JSON `json:"video" gorm:"video"`
	Cover datatypes.JSON `json:"cover" gorm:"cover"`
}

type BarragesList []Barrage

type BarragesJson [][]interface{}

type GetVideoBarrageListReceiveStruct struct {
	ID string `json:"id"`
}

type barrageInfo struct {
	Time     int       `json:"time"`
	Text     string    `json:"text"`
	SendTime time.Time `json:"sendTime"`
}

type BarrageInfoList []barrageInfo

type GetVideoCommentReceiveStruct struct {
	PageInfo common.PageInfo `json:"pageInfo"`
	VideoID  uint            `json:"video_id" binding:"required"`
}

// 评论信息
type commentsInfo struct {
	ID              uint             `json:"id"`
	CommentID       uint             `json:"comment_id"`
	CommentFirstID  uint             `json:"comment_first_id"`
	CreatedAt       time.Time        `json:"created_at"`
	Context         string           `json:"context"`
	Uid             uint             `json:"uid"`
	Username        string           `json:"username"`
	Photo           string           `json:"photo"`
	CommentUserID   uint             `json:"comment_user_id"`
	CommentUserName string           `json:"comment_user_name"`
	LowerComments   commentsInfoList `json:"lowerComments"`
}

type commentsInfoList []*commentsInfo

type GetArticleContributionCommentsResponseStruct struct {
	Id             uint             `json:"id"`
	Comments       commentsInfoList `json:"comments"`
	CommentsNumber int              `json:"comments_number"`
}

type GetVideoContributionByIDReceiveStruct struct {
	VideoID uint `json:"video_id"`
}

// Info 视频信息
type Info struct {
	ID             uint             `json:"id"`
	Uid            uint             `json:"uid" `
	Title          string           `json:"title" `
	Video          string           `json:"video"`
	Video720p      string           `json:"video_720p"`
	Video480p      string           `json:"video_480p"`
	Video360p      string           `json:"video_360p"`
	Cover          string           `json:"cover" `
	VideoDuration  int64            `json:"video_duration"`
	Label          []string         `json:"label"`
	Introduce      string           `json:"introduce"`
	Heat           int              `json:"heat"`
	BarrageNumber  int              `json:"barrageNumber"`
	Comments       commentsInfoList `json:"comments"`
	IsLike         bool             `json:"is_like"`
	IsCollect      bool             `json:"is_collect"`
	CommentsNumber int              `json:"comments_number"`
	CreatorInfo    creatorInfo      `json:"creatorInfo"`
	CreatedAt      time.Time        `json:"created_at"`
}

// 创作者信息
type creatorInfo struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Signature   string `json:"signature"`
	IsAttention bool   `json:"is_attention"`
}

// 推荐视频信息
type recommendVideo struct {
	ID            uint      `json:"id"`
	Uid           uint      `json:"uid" `
	Title         string    `json:"title" `
	Video         string    `json:"video"`
	Cover         string    `json:"cover" `
	VideoDuration int64     `json:"video_duration"`
	Label         []string  `json:"label"`
	Introduce     string    `json:"introduce"`
	Heat          int       `json:"heat"`
	BarrageNumber int       `json:"barrageNumber"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
}
type RecommendList []recommendVideo

type Response struct {
	VideoInfo     Info          `json:"videoInfo"`
	RecommendList RecommendList `json:"recommendList"`
}

type SendVideoBarrageReceiveStruct struct {
	Author string  `json:"author"`
	Color  uint    `json:"color" binding:"required"`
	ID     string  `json:"id" binding:"required"`
	Text   string  `json:"text" binding:"required"`
	Time   float64 `json:"time"`
	Type   uint    `json:"type"`
	Token  string  `json:"token" binding:"required"`
}
