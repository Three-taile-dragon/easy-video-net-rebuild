package model

import (
	"dragonsss.cn/evn_common/calculate"
	"dragonsss.cn/evn_common/conversion"
	comments2 "dragonsss.cn/evn_common/model/article/comments"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_common/model/video/barrage"
	"dragonsss.cn/evn_common/model/video/comments"
	"time"
)

type GetDiscussVideoListItem struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Photo     string    `json:"photo"`
	Comment   string    `json:"comment"`
	Cover     string    `json:"cover"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type GetDiscussVideoListStruct []GetDiscussVideoListItem

func GetDiscussVideoListResponse(cml *comments.CommentList, localhost string, tencentOssHost string) interface{} {
	list := make(GetDiscussVideoListStruct, 0)
	for _, v := range *cml {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo, localhost, tencentOssHost)
		cover, _ := conversion.FormattingJsonSrc(v.VideoInfo.Cover, localhost, tencentOssHost)
		list = append(list, GetDiscussVideoListItem{
			ID:        v.ID,
			Username:  v.UserInfo.Username,
			Photo:     photo,
			Comment:   v.Context,
			Cover:     cover,
			Title:     v.VideoInfo.Title,
			CreatedAt: v.CreatedAt,
		})
	}
	return list
}

type GetDiscussArticleListItem struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Photo     string    `json:"photo"`
	Comment   string    `json:"comment"`
	Cover     string    `json:"cover"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type GetDiscussArticleListStruct []GetDiscussArticleListItem

func GetDiscussArticleListResponse(cml *comments2.CommentList, localhost string, tencentOssHost string) interface{} {
	list := make(GetDiscussArticleListStruct, 0)
	for _, v := range *cml {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo, localhost, tencentOssHost)
		cover, _ := conversion.FormattingJsonSrc(v.ArticleInfo.Cover, localhost, tencentOssHost)
		list = append(list, GetDiscussArticleListItem{
			ID:        v.ID,
			Username:  v.UserInfo.Username,
			Photo:     photo,
			Comment:   v.Context,
			Cover:     cover,
			Title:     v.ArticleInfo.Title,
			CreatedAt: v.CreatedAt,
		})
	}
	return list
}

type GetDiscussBarrageListItem struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Photo     string    `json:"photo"`
	Comment   string    `json:"comment"`
	Cover     string    `json:"cover"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type GetDiscussBarrageListStruct []GetDiscussBarrageListItem

func GetDiscussBarrageListResponse(cml *barrage.BarragesList, localhost string, tencentOssHost string) interface{} {
	list := make(GetDiscussBarrageListStruct, 0)
	for _, v := range *cml {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo, localhost, tencentOssHost)
		cover, _ := conversion.FormattingJsonSrc(v.VideoInfo.Cover, localhost, tencentOssHost)
		list = append(list, GetDiscussBarrageListItem{
			ID:        v.ID,
			Username:  v.UserInfo.Username,
			Photo:     photo,
			Comment:   v.Text,
			Cover:     cover,
			Title:     v.VideoInfo.Title,
			CreatedAt: v.CreatedAt,
		})
	}
	return list
}

func SearchVideoResponse(videoList *video.VideosContributionList, localhost string, tencentOssHost string) (interface{}, error) {
	//处理视频
	vl := make(videoInfoList, 0)
	for _, lk := range *videoList {
		cover, _ := conversion.FormattingJsonSrc(lk.Cover, localhost, tencentOssHost)
		videoSrc, _ := conversion.FormattingJsonSrc(lk.Video, localhost, tencentOssHost)
		info := VideoInfo{
			ID:            lk.ID,
			Uid:           lk.Uid,
			Title:         lk.Title,
			Video:         videoSrc,
			Cover:         cover,
			VideoDuration: lk.VideoDuration,
			Label:         conversion.StringConversionMap(lk.Label),
			Introduce:     lk.Introduce,
			Heat:          lk.Heat,
			BarrageNumber: len(lk.Barrage),
			Username:      lk.UserInfo.Username,
			CreatedAt:     lk.CreatedAt,
		}
		vl = append(vl, info)
	}
	return vl, nil
}

type UserInfo struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Photo       string `json:"photo"`
	Signature   string `json:"signature"`
	IsAttention bool   `json:"is_attention"`
}

type UserInfoList []UserInfo

func SearchUserResponse(userList *user.UserList, aids []uint, localhost string, tencentOssHost string) (interface{}, error) {
	list := make(UserInfoList, 0)
	for _, v := range *userList {
		photo, _ := conversion.FormattingJsonSrc(v.Photo, localhost, tencentOssHost)
		list = append(list, UserInfo{
			ID:          v.ID,
			Username:    v.Username,
			Photo:       photo,
			Signature:   v.Signature,
			IsAttention: calculate.ArrayIsContain(aids, v.ID),
		})
	}
	return list, nil
}
